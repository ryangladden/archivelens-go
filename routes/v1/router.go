package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/ryangladden/archivelens-go/handler"
	// "github.com/ryangladden/archivelens-go/handlers/middleware"
)

type Router struct {
	userHandler     *handler.UserHandler
	authHandler     *handler.AuthHandler
	documentHandler *handler.DocumentHandler
	personHandler   *handler.PersonHandler
	routes          *gin.Engine
}

func NewRouter(userHandler *handler.UserHandler, authHandler *handler.AuthHandler, documentHandler *handler.DocumentHandler, personHandler *handler.PersonHandler) *Router {
	r := gin.Default()

	router := &Router{
		userHandler:     userHandler,
		authHandler:     authHandler,
		documentHandler: documentHandler,
		personHandler:   personHandler,
		routes:          r,
	}

	router.registerRoutes()
	return router
}

func (r *Router) Run(addr string) error {
	return r.routes.Run(addr)
}

func (r *Router) registerRoutes() {
	r.routes.GET("/", func(c *gin.Context) { c.JSON(200, gin.H{"message": "pong"}) })

	api := r.routes.Group("/api")
	v1 := api.Group("/v1")
	users := v1.Group("/users")
	{
		users.POST("", r.userHandler.CreateUser)
		// users.GET("me", r.authHandler.AuthenticateMiddleware(), r.userHandler.GetMe)
		// 	users.PUT("", CreateUser)
		// 	users.PATCH("", UpdateUser)
		// 	users.DELETE("", DeleteUser)
	}
	auth := v1.Group("/auth")
	{
		auth.POST("/login", r.authHandler.CreateAuth)
		auth.DELETE("/logout", r.authHandler.DeleteAuth)
		auth.GET("/me", r.authHandler.AuthenticateMiddleware(), r.authHandler.GetSession)
	}
	documents := v1.Group("/documents")
	documents.Use(r.authHandler.AuthenticateMiddleware())
	{
		// documents.GET("/:id", r.documentHandler.GetDocument)
		documents.POST("/upload", r.documentHandler.CreateDocument)
		// 	documents.GET("/:id", GetDocument)
		// 	documents.PATCH("/:id", UpdateDocument)
		// 	documents.DELETE("/:id", DeleteDocument)
	}
	persons := v1.Group("/persons")
	persons.Use(r.authHandler.AuthenticateMiddleware())
	{
		// persons.GET("", GetPersons)
		persons.POST("", r.personHandler.CreatePerson)
		// 	persons.GET("/:id", GetPerson)
		// 	persons.PATCH("/:id", UpdatePerson)
		// 	persons.DELETE("/:id", DeletePerson)
	}
}
