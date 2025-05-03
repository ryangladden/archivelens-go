package handler

import (
	"github.com/rs/zerolog/log"

	"github.com/gin-gonic/gin"
	errs "github.com/ryangladden/archivelens-go/err"
	requests "github.com/ryangladden/archivelens-go/request"
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

func (h *UserHandler) CreateUser(c *gin.Context) {
	var createUserRequest requests.CreateUserRequest
	log.Info().Msg("POST /api/v1/users")
	if err := c.BindJSON(&createUserRequest); err != nil {
		log.Error().Err(err).Msg("Invalid request body for creating user")
		c.JSON(400, gin.H{"error": "invalid request body"})
		return
	}
	if err := h.userService.CreateUser(&createUserRequest); err != nil {
		if err == errs.ErrConflict {
			c.JSON(409, gin.H{"error": "user with this email already exists"})
			return
		} else {
			c.JSON(500, gin.H{"error": "internal server error"})
			return
		}
	}
	c.JSON(201, gin.H{"message": "user created successfully"})
}
