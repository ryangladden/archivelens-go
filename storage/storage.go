package storage

// import (
// 	"context"
// 	"fmt"
// 	"net/url"

// 	"github.com/rs/zerolog/log"
// 	"github.com/aws/aws-sdk-go-v2/config"
// 	"github.com/aws/aws-sdk-go-v2/service/s3"
// 	"github.com/aws/smithy-go/endpoints"
// )

// type S3Manager struct {
// 	bucketName string
// 	service *s3.S3
// 	config *aws.Config
// }

// func NewS3Manager(s3Endpoint string, s3BucketName string, s3Location string) {

// 	endpointURL, err := url.Parse(s3Endpoint)
// 	if err != nil {
// 		log.Fatal().Err(err).Msgf("Failed to parse S3 endpoint. AWS_ENDPOINT=%s", s3Endpoint)
// 	}
// 	disableSSL := true
// 	config.NewEnvConfig()
// 	session := session.Must(session.NewSession())
// 	service := s3.New(sess)
// 	ctx := context.Background()
// }
