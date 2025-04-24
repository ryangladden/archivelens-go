package handlers

import (
	"fmt"
	"io"

	// "github.com/ryangladden/archivelens-go/model"
	"github.com/gin-gonic/gin"
	"github.com/ryangladden/archivelens-go/requests"
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

func (h *UserHandler) CreateUser(c *gin.Context) {
	var createUserRequest requests.CreateUserRequest

	if err := c.BindJSON(&createUserRequest); err != nil {
		c.JSON(400, gin.H{"error": "invalid request body"})
		return
	}
	fmt.Sprintf("Name: %s\nPassword: %s\nEmail: %s\n", createUserRequest.Name, createUserRequest.Password, createUserRequest.Email)
	h.userService.CreateUser(&createUserRequest)
	c.JSON(201, gin.H{"message": "user created successfully"})
}
