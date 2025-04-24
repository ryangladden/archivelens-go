package handlers

import (
	"fmt"
	"io"

	// "github.com/ryangladden/archivelens-go/model"
	// "github.com/ryangladden/archivelens-go/requests"
	"github.com/gin-gonic/gin"
	"github.com/ryangladden/archivelens-go/service"
)

type UserHandler struct {
	userService *service.UserService
}

func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

func CreateUserHandler(c *gin.Context) {
	body, err := io.ReadAll(c.Request.Body)
	fmt.Println((body))
	if err != nil {
		c.JSON(500, gin.H{"error": "failed to read request body"})
	}
	c.JSON(200, gin.H{"message": "pong"})
}
