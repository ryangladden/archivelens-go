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
	userHandler       *handlers.UserHandler
	userService       *service.UserService
	userDao           *db.UserDAO
	router            *routes.Router
}

func NewServer() *Server {
	connectionManager := db.NewConnectionManager("localhost", 5432, "postgres", "postgres", "archive-lens-dev")
	userDao := db.NewUserDAO(connectionManager)
	userService := service.NewUserService(userDao)
	userHandler := handlers.NewUserHandler(userService)
	router := routes.NewRouter(userHandler)

	connectionManager.Connect()

	return &Server{
		userHandler: userHandler,
		userService: userService,
		userDao:     userDao,
		router:      router,
	}
}

func main() {

	zerolog.TimeFieldFormat = "2006-01-02 15:04:05"
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout}).With().Timestamp().Logger()

	server := NewServer()
	server.router.Run(":8080")
}
