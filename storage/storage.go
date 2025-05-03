package storage

import (
	"context"
	"fmt"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func CreateBucket() {
	ctx := context.Background()
	endpoint := "localhost:9000"
	accessKeyID := "9tAeZxTVqaHMWAz0nbAh"
	secretAccessKey := "nFLdSPVcsrdeGxtSFmxNea8V5XtPPfk32JzeiOIQ"
	useSSL := false

	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
	})

	if err != nil {
		fmt.Printf("Error: %v", err)
	}

	bucketName := "archive-lens"
	location := "us-east-1"

	err = minioClient.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{
		Region: location,
	})

	if err != nil {
		exists, errBucketExists := minioClient.BucketExists(ctx, bucketName)
		if errBucketExists == nil && exists {
			fmt.Println("bucket exists dude")
		} else {
			fmt.Printf("Error: %v", err)
		}
	} else {
		fmt.Println("Bucket created dawg")
	}
}
