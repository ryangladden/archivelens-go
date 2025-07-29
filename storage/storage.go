package storage

import (
	"context"
	"errors"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

type StorageManager struct {
	Client     *s3.Client
	bucketName string
	presigner  *s3.PresignClient
}

func NewStorageManager(s3Endpoint string, s3BucketName string, s3Location string) *StorageManager {
	// id := "9fA1jgAKDsQVrdkMExwx"
	// key := "JGZhXCL1qQZy8KrDZgiMc3UmOKNrN1yRO2twyyGI"
	id := os.Getenv("AWS_ACCESS_KEY_ID")
	key := os.Getenv("AWS_SECRET_ACCESS_KEY")
	options := s3.Options{
		Region:       s3Location,
		BaseEndpoint: &s3Endpoint,
		UsePathStyle: true,
		Credentials: aws.CredentialsProviderFunc(func(ctx context.Context) (aws.Credentials, error) {
			return aws.Credentials{
				AccessKeyID:     id,
				SecretAccessKey: key,
			}, nil
		}),
	}

	client := s3.New(options)
	sm := &StorageManager{
		Client:     client,
		bucketName: s3BucketName,
		presigner:  s3.NewPresignClient(client),
	}

	sm.createBucket(sm.bucketName)
	return sm
}

func (s *StorageManager) createBucket(bucket string) {
	log.Debug().Msgf("Creating bucket %s", bucket)
	input := s3.CreateBucketInput{
		Bucket: &bucket,
	}
	_, err := s.Client.CreateBucket(context.Background(), &input)
	if err != nil {
		var owned *types.BucketAlreadyOwnedByYou
		var exists *types.BucketAlreadyExists
		if errors.As(err, &owned) {
			log.Info().Msgf("You already own bucket %s", bucket)
		} else if errors.As(err, &exists) {
			log.Info().Msgf("Bucket %s already exists", bucket)
		} else {
			log.Fatal().Err(err).Msgf("Failed to create bucket %s. Initialization failed", bucket)
		}
	}
}

func (s *StorageManager) UploadFile(file *multipart.FileHeader, key string) error {

	reader, err := file.Open()
	if err != nil {
		log.Error().Err(err).Msgf("Error opening file %s", file.Filename)
		return err
	}
	contentType := file.Header.Get("Content-Type")
	contentDisposition := "inline"
	input := s3.PutObjectInput{
		Bucket:             &s.bucketName,
		Key:                &key,
		Body:               reader,
		ContentType:        &contentType,
		ContentDisposition: &contentDisposition,
	}
	_, err = s.Client.PutObject(context.Background(), &input)
	if err != nil {
		log.Error().Err(err).Msgf("Error uploading file %s with key %s to bucket %s", file.Filename, key, s.bucketName)
		return err
	}
	return nil
}

func (s *StorageManager) GeneratePresignedURL(key *string) *string {
	if key != nil {
		input := s3.GetObjectInput{
			Bucket: &s.bucketName,
			Key:    key,
		}
		request, err := s.presigner.PresignGetObject(context.Background(), &input, s3.WithPresignExpires(time.Duration(15*int(time.Second))))
		if err != nil {
			log.Error().Err(err).Msgf("Failed getting presigned URL for object with key %s", *key)
		}
		return &request.URL
	}
	return nil
}

func GenerateObjectKey(base string, id uuid.UUID, newFileName string, filename string) *string {
	extension := strings.ToLower(filepath.Ext(filename))
	if extension == ".jpeg" {
		extension = ".jpg"
	}
	key := filepath.Join(base, id.String(), newFileName+extension)
	return &key
}
