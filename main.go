package main

import (
	"fmt"
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

	thing := []string{"joe", "mama"}
	fmt.Println(thing)
	zerolog.TimeFieldFormat = "2006-01-02 15:04:05"
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout}).With().Timestamp().Logger()

	server := server.NewServer()
	server.Run(":8080")
}
