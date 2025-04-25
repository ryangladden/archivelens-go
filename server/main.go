package main

import (
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/ryangladden/archivelens-go/db"
	"github.com/ryangladden/archivelens-go/handlers"
	"github.com/ryangladden/archivelens-go/routes/v1"
	"github.com/ryangladden/archivelens-go/service"
)

var (
	LogJSON = zerolog.New(zerolog.ConsoleWriter{Out: os.Stdout}).With().Timestamp().Logger()
)

type Server struct {
	connectionManager *db.ConnectionManager

	userHandler *handlers.UserHandler
	authHandler *handlers.AuthHandler

	userService *service.UserService
	authService *service.AuthService

	userDao *db.UserDAO
	authDao *db.AuthDAO

	router *routes.Router
}

func NewServer() *Server {
	connectionManager := db.NewConnectionManager("localhost", 5432, "postgres", "postgres", "archive-lens-dev")

	userDao := db.NewUserDAO(connectionManager)
	userService := service.NewUserService(userDao)
	userHandler := handlers.NewUserHandler(userService)

	authDao := db.NewAuthDAO(connectionManager)
	authService := service.NewAuthService(authDao, userDao)
	authHandler := handlers.NewAuthHandler(authService)

	router := routes.NewRouter(userHandler, authHandler)

	connectionManager.Connect()
	err := db.SetUp(connectionManager)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to set up database")
	}

	return &Server{
		connectionManager: connectionManager,

		userHandler: userHandler,
		authHandler: authHandler,

		userService: userService,
		authService: authService,

		userDao: userDao,
		authDao: authDao,

		router: router,
	}
}

func main() {

	zerolog.TimeFieldFormat = "2006-01-02 15:04:05"
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout}).With().Timestamp().Logger()

	server := NewServer()
	server.router.Run(":8080")
}
