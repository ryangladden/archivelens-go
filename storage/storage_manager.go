package storage

import (
	"context"
	"mime/multipart"
	"net/url"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"

	// "github.com/klippa-app/go-pdfium"
	// "github.com/klippa-app/go-pdfium/requests"
	errs "github.com/ryangladden/archivelens-go/err"

	// "path/filepath"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type StorageManager struct {
	minioClient *minio.Client
	bucketName  string
}

func NewStorageManager(s3Endpoint string, s3AccessKeyId string, s3SecretAccessKey string, s3BucketName string, s3Location string) *StorageManager {
	minioClient := createClient(s3Endpoint, s3AccessKeyId, s3SecretAccessKey)
	s3Init(minioClient, s3BucketName, s3Location)

	return &StorageManager{
		minioClient: minioClient,
		bucketName:  s3BucketName,
	}
}

func s3Init(minioClient *minio.Client, s3bucketName string, s3location string) {
	ctx := context.Background()
	err := minioClient.MakeBucket(ctx, s3bucketName, minio.MakeBucketOptions{
		Region: s3location,
	})

	if err != nil {
		exists, errBucketExists := minioClient.BucketExists(ctx, s3bucketName)
		if errBucketExists == nil && exists {
			log.Info().Msgf("Bucket \"%s\" exists, skipping bucket creation", s3bucketName)
		} else {
			log.Fatal().Err(err).Msgf("Failed to initialize bucket \"%s\"", s3bucketName)
		}
	}
}

func (sm *StorageManager) GenerateObjectKey(filename string, id uuid.UUID, directory string) *string {
	extension := strings.ToLower(filepath.Ext(filename))
	key := filepath.Join(directory, id.String()+extension)
	return &key
}

func (sm *StorageManager) UploadFile(file *multipart.FileHeader, key *string) error {

	if sm.minioClient == nil {
		log.Error().Msg("Minio client is offline")
		return errs.ErrStorage
	}
	log.Debug().Msgf("Content-Type: %s", file.Header.Get("Content-Type"))
	reader, err := file.Open()
	if err != nil {
		log.Error().Err(err).Msgf("Error opening file: %s", file.Filename)
		return errs.ErrStorage
	}
	defer reader.Close()

	ctx := context.Background()
	_, err = sm.minioClient.PutObject(ctx, sm.bucketName, *key, reader, file.Size, minio.PutObjectOptions{ContentType: file.Header.Get("Content-Type"), ContentDisposition: "inline"})
	if err != nil {
		log.Error().Err(err).Msg("Failed to put in the bucket")
	}
	return nil
}

func (sm *StorageManager) GeneratePresignedURL(key *string) *string {
	if key != nil {
		reqParams := make(url.Values)
		url, err := sm.minioClient.PresignedGetObject(context.Background(), sm.bucketName, *key, 15*time.Second, reqParams)
		if err == nil {
			presignedUrl := url.String()
			return &presignedUrl
		}
	}
	return nil
}

func createClient(s3Endpoint string, s3AccessKeyId string, s3SecretAccessKey string) *minio.Client {
	minioClient, err := minio.New(s3Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(s3AccessKeyId, s3SecretAccessKey, ""),
		Secure: false,
	})
	if err != nil {
		log.Fatal().Err(err).Msgf("Unable to create S3 client at endpoint %s", s3Endpoint)
	}
	log.Info().Msgf("New S3 client accessing %s", s3Endpoint)
	return minioClient
}

// func generateThumbnail(file *multipart.FileHeader) {
// 	reader, err := file.Open()

// 	doc, err := pdfium.Pdfium.OpenDocument(
// 		&requests.OpenDocument{
// 			FileReader: reader,
// 		},
// 	)
// 	img, err := pdfium.Pdfium.RenderPageInPixels(
// 		&requests.RenderPageInPixels{
// 			Page: requests.Page{
// 				ByIndex: &requests.PageByIndex{
// 					Document: doc,
// 					Index:    0,
// 			},
// 			Width: 1,
// 			Height: 1,
// 		},
// 	)
// 	doc.

// }
