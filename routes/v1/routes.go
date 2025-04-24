package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/ryangladden/archivelens-go/handlers"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()
	router.GET("/", func(c *gin.Context) { c.JSON(200, gin.H{"message": "pong"}) })

	api := router.Group("/api")
	v1 := api.Group("/v1")
	users := v1.Group("/users")
	{
		users.GET("", handlers.CreateUserHandler)
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
	return router
}
