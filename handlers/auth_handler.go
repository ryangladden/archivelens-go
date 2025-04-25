package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"

	"github.com/ryangladden/archivelens-go/errs"
	"github.com/ryangladden/archivelens-go/requests"
	"github.com/ryangladden/archivelens-go/service"
)

type AuthHandler struct {
	authService *service.AuthService
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

func (h *AuthHandler) CreateAuth(c *gin.Context) {
	var loginRequest requests.LoginRequest
	if err := c.BindJSON(&loginRequest); err != nil {
		log.Error().Err(err).Msg("Invalid request body for logging in")
		c.JSON(400, gin.H{"error": "invalid request body"})
	}

	authToken, err := h.authService.CreateAuth(loginRequest)
	if err != nil {
		if err == errs.ErrNotFound {
			c.JSON(404, gin.H{"error": "user not found"})
		} else if err == errs.ErrUnauthorized {
			c.JSON(401, gin.H{"error": "unauthorized access"})
		}
	}

	c.JSON(200, gin.H{"auth_token": authToken})
}
