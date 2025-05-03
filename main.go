package main

import (
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/ryangladden/archivelens-go/server"
)

var (
	s3Endpoint        string
	s3AccessKeyId     string
	s3SecretAccessKey string
	s3BucketName      string
	s3Location        string
)

func main() {

	zerolog.TimeFieldFormat = "2006-01-02 15:04:05"
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout}).With().Timestamp().Logger()

	server := server.NewServer()
	server.Init(":8080")
}
