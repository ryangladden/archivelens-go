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

	s3Endpoint   string
	s3BucketName string
	s3Location   string

	valkeyEndpoint string
)

type Server struct {
	connectionManager *db.ConnectionManager
	storageManager    *storage.StorageManager

	// userHandler     *handler.UserHandler
	authHandler     *handler.AuthHandler
	documentHandler *handler.DocumentHandler
	personHandler   *handler.PersonHandler

	authService     *service.AuthService
	documentService *service.DocumentService
	personService   *service.PersonService

	// userDao     *db.UserDAO
	authDao     *db.AuthDAO
	documentDao *db.DocumentDAO
	personDao   *db.PersonDAO

	router *routes.Router
}

func NewServer() *Server {
	getEnvironmentVariables()

	connectionManager := db.NewConnectionManager(postgresHost, postgresPort, postgresUsername, postgresPassword, postgresDb)
	// storageManager := storage.NewStorageManager(s3Endpoint, s3AccessKeyId, s3SecretAccessKey, s3BucketName, s3Location)
	storageManager := storage.NewStorageManager(s3Endpoint, s3BucketName, s3Location)

	// userDao := db.NewUserDAO(connectionManager)
	// userService := service.NewUserService(userDao)
	// userHandler := handler.NewUserHandler(userService)

	authDao := db.NewAuthDAO(connectionManager)
	authService := service.NewAuthService(authDao)
	authHandler := handler.NewAuthHandler(authService)

	documentDao := db.NewDocumentDAO(connectionManager)
	documentService := service.NewDocumentService(documentDao, storageManager)
	documentHandler := handler.NewDocumentHandler(documentService)

	personDao := db.NewPersonDAO(connectionManager)
	personService := service.NewPersonService(personDao, storageManager)
	personHandler := handler.NewPersonHandler(personService)

	router := routes.NewRouter(authHandler, documentHandler, personHandler)

	return &Server{
		connectionManager: connectionManager,
		storageManager:    storageManager,

		// userHandler:     userHandler,
		authHandler:     authHandler,
		documentHandler: documentHandler,
		personHandler:   personHandler,

		// userService:     userService,
		authService:     authService,
		documentService: documentService,
		personService:   personService,

		// userDao:     userDao,
		authDao:     authDao,
		documentDao: documentDao,
		personDao:   personDao,

		router: router,
	}
}

func (s *Server) Run(hostname string) {
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

	s3Endpoint = os.Getenv("AWS_ENDPOINT")
	// s3AccessKeyId = os.Getenv("S3_ACCESS_KEY_ID")
	// s3SecretAccessKey = os.Getenv("S3_SECRET_ACCESS_KEY")
	s3BucketName = os.Getenv("AWS_BUCKET_NAME")
	s3Location = os.Getenv("AWS_REGION")
}
