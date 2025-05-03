package storage

import (
	"context"
	"mime/multipart"

	"github.com/rs/zerolog/log"
	errs "github.com/ryangladden/archivelens-go/err"

	// "path/filepath"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type StorageManager struct {
	minioClient *minio.Client
}

func NewStorageManager(s3Endpoint string, s3AccessKeyId string, s3SecretAccessKey string, s3BucketName string, s3Location string) *StorageManager {
	minioClient, err := createClient(s3Endpoint, s3AccessKeyId, s3SecretAccessKey)
	if err != nil {
		panic(err)
	}
	return &StorageManager{
		minioClient: minioClient,
	}
}

func (sm *StorageManager) S3Init(s3bucketName string, s3location string) error {
	ctx := context.Background()
	err := sm.minioClient.MakeBucket(ctx, s3bucketName, minio.MakeBucketOptions{
		Region: s3location,
	})

	if err != nil {
		exists, errBucketExists := sm.minioClient.BucketExists(ctx, s3bucketName)
		if errBucketExists == nil && exists {
			log.Info().Msgf("Bucket \"%s\" exists, skipping bucket creation", s3bucketName)
		} else {
			log.Error().Err(err).Msgf("Failed to initialize bucket \"%s\"", s3bucketName)
			panic(err)
		}
	}
	return nil
}

func (sm *StorageManager) UploadFile(file *multipart.FileHeader, bucketName string, key string) error {
	ctx := context.Background()
	reader, err := file.Open()
	if sm.minioClient == nil {
		panic(err)
	}
	if err != nil {
		log.Error().Err(err).Msgf("Error uploading file: %s", file.Filename)
		return errs.ErrStorage
	}
	defer reader.Close()
	sm.minioClient.PutObject(ctx, bucketName, key, reader, file.Size, minio.PutObjectOptions{})
	return nil
}

func createClient(s3Endpoint string, s3AccessKeyId string, s3SecretAccessKey string) (*minio.Client, error) {
	minioClient, err := minio.New(s3Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(s3AccessKeyId, s3SecretAccessKey, ""),
		Secure: false,
	})
	if err != nil {
		log.Error().Err(err).Msgf("Error connection to S3 instance at endpoint %s", s3Endpoint)
		return nil, errs.ErrStorage
	}
	log.Info().Msgf("New S3 client accessing %s", s3Endpoint)
	return minioClient, nil
}
