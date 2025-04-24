package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/ryangladden/archivelens-go/handlers"
)

type Router struct {
	userHandler *handlers.UserHandler
	routes      *gin.Engine
}

func NewRouter(userHandler *handlers.UserHandler) *Router {
	r := gin.Default()

	router := &Router{
		userHandler: userHandler,
		routes:      r,
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
		// 	users.PUT("", CreateUser)
		// 	users.PATCH("", UpdateUser)
		// 	users.DELETE("", DeleteUser)
	}
	// documents := v1.Group("/documents")
	// {
	// 	documents.GET("", GetDocuments)
	// 	documents.POST("", CreateDocument)
	// 	documents.GET("/:id", GetDocument)
	// 	documents.PATCH("/:id", UpdateDocument)
	// 	documents.DELETE("/:id", DeleteDocument)
	// }
	// persons := v1.Group("/persons")
	// {
	// 	persons.GET("", GetPersons)
	// 	persons.POST("", CreatePerson)
	// 	persons.GET("/:id", GetPerson)
	// 	persons.PATCH("/:id", UpdatePerson)
	// 	persons.DELETE("/:id", DeletePerson)
	// }
}
