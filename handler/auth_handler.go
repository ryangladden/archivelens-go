package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"

	errs "github.com/ryangladden/archivelens-go/err"
	"github.com/ryangladden/archivelens-go/model"
	"github.com/ryangladden/archivelens-go/request"
	"github.com/ryangladden/archivelens-go/response"
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
	var loginRequest request.LoginRequest
	if err := c.BindJSON(&loginRequest); err != nil {
		log.Error().Err(err).Msg("Invalid request body for logging in")
		c.JSON(400, gin.H{"error": "invalid request body"})
		return
	}

	authToken, user, err := h.authService.CreateAuth(loginRequest)

	if err != nil {
		if err == errs.ErrNotFound {
			c.JSON(404, gin.H{"error": "user not found"})
			return
		} else if err == errs.ErrUnauthorized {
			c.JSON(401, gin.H{"error": "unauthorized access"})
			return
		}
	}

	http.SetCookie(c.Writer, createCookie(authToken))

	c.JSON(200, user)
}

func (h *AuthHandler) DeleteAuth(c *gin.Context) {
	authToken, err := c.Cookie("archive_lens_access_token")
	if err != nil {
		c.JSON(204, nil)
		return
	}
	h.authService.DeleteAuth(authToken)
	c.SetCookie("archive_lens_access_token", "", -1, "/", "", false, true)
	c.JSON(200, nil)
}

func (h *AuthHandler) GetSession(c *gin.Context) {
	user := getUserFromContext(c)
	if user != nil {
		var response = &response.LoginResponse{
			Email:     user.Email,
			FirstName: user.FirstName,
			LastName:  user.LastName,
		}
		c.JSON(200, response)
		return
	}
	c.AbortWithStatus(500)
}

func (h *AuthHandler) AuthenticateMiddleware() gin.HandlerFunc {
	log.Debug().Msg("AuthenticatedMiddleware implemented")
	return func(c *gin.Context) {
		log.Debug().Msg("AuthenticatedMiddleware called")
		token, err := c.Cookie("archive_lens_access_token")
		if err != nil {
			c.JSON(401, gin.H{"error": "Unauthorized"})
			return
		}
		log.Debug().Msgf("Token from database: %s", token)
		user, err := h.authService.ValidateToken(token)

		if err == errs.ErrDB {
			c.AbortWithStatus(500)
			return
		} else if err == errs.ErrUnauthorized {
			c.AbortWithStatus(401)
			return
		}
		log.Debug().Msgf("Setting user in gin context: %s", user.ID)
		c.Set("user", user.ID)
		c.Next()
	}
}

func createCookie(token string) *http.Cookie {
	return &http.Cookie{
		Name:     "archive_lens_access_token",
		Value:    token,
		Path:     "/",
		Domain:   "",
		Expires:  time.Now().Add(24 * time.Hour * 180),
		MaxAge:   86400 * 180,
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteStrictMode,
	}
}

func getUserFromContext(c *gin.Context) *model.User {
	user, exists := c.Get("user")
	if exists {
		if u, ok := user.(*model.User); ok {
			return u
		}
	}
	return nil
}
