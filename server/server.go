package server

import (
	"os"
	"strconv"

	"github.com/rs/zerolog"

	"github.com/ryangladden/archivelens-go/db"
	"github.com/ryangladden/archivelens-go/handler"
	"github.com/ryangladden/archivelens-go/routes/v1"
	"github.com/ryangladden/archivelens-go/service"
	"github.com/ryangladden/archivelens-go/storage"
)

var (
	LogJSON = zerolog.New(zerolog.ConsoleWriter{Out: os.Stdout}).With().Timestamp().Logger()

	postgresHost     string
	postgresPort     int
	postgresUsername string
	postgresPassword string
	postgresDb       string

	s3Endpoint        string
	s3AccessKeyId     string
	s3SecretAccessKey string
	s3BucketName      string
	s3Location        string
)

type Server struct {
	connectionManager *db.ConnectionManager
	storageManager    *storage.StorageManager

	userHandler *handler.UserHandler
	authHandler *handler.AuthHandler

	userService *service.UserService
	authService *service.AuthService

	userDao *db.UserDAO
	authDao *db.AuthDAO

	router *routes.Router
}

func NewServer() *Server {
	getEnvironmentVariables()

	connectionManager := db.NewConnectionManager(postgresHost, postgresPort, postgresUsername, postgresPassword, postgresDb)
	storageManager := storage.NewStorageManager(s3Endpoint, s3AccessKeyId, s3SecretAccessKey, s3BucketName, s3Location)

	userDao := db.NewUserDAO(connectionManager)
	userService := service.NewUserService(userDao)
	userHandler := handler.NewUserHandler(userService)

	authDao := db.NewAuthDAO(connectionManager)
	authService := service.NewAuthService(authDao, userDao)
	authHandler := handler.NewAuthHandler(authService)

	documentDao := db.NewDocumentDAO(connectionManager)
	documentService := service.NewDocumentService(documentDao)
	documentHandler := handler.NewDocumentHandler(documentService)

	router := routes.NewRouter(userHandler, authHandler, documentHandler)

	return &Server{
		connectionManager: connectionManager,
		storageManager:    storageManager,

		userHandler: userHandler,
		authHandler: authHandler,

		userService: userService,
		authService: authService,

		userDao: userDao,
		authDao: authDao,

		router: router,
	}
}

func (s *Server) Init(hostname string) {
	s.connectionManager.Init()
	s.storageManager.S3Init(s3BucketName, s3Location)
	s.router.Run(hostname)
}

func getEnvironmentVariables() {
	var err error
	postgresHost = os.Getenv("POSTGRES_HOST")
	postgresPort, err = strconv.Atoi(os.Getenv("POSTGRES_PORT"))
	if err != nil {
		panic(err)
	}
	postgresUsername = os.Getenv("POSTGRES_USERNAME")
	postgresPassword = os.Getenv("POSTGRES_PASSWORD")
	postgresDb = os.Getenv("POSTGRES_DB")

	s3Endpoint = os.Getenv("S3_ENDPOINT")
	s3AccessKeyId = os.Getenv("S3_ACCESS_KEY_ID")
	s3SecretAccessKey = os.Getenv("S3_SECRET_ACCESS_KEY")
	s3BucketName = os.Getenv("S3_BUCKET_NAME")
	s3Location = os.Getenv("S3_LOCATION")
}
