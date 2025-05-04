package middleware

import (
	"github.com/gin-gonic/gin"
	errs "github.com/ryangladden/archivelens-go/err"
	"github.com/ryangladden/archivelens-go/service"
)

func AuthenticateMiddleware(authService *service.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := c.Cookie("archive_lens_access_token")
		if err != nil {
			c.JSON(401, gin.H{"error": "Unauthorized"})
		}

		user, err := authService.ValidateToken(token)
		if err == errs.ErrDB {
			c.AbortWithStatus(500)
			return
		} else if err == errs.ErrUnauthorized {
			c.AbortWithStatus(401)
			return
		}

		c.Set("user", user.ID)
		c.Next()
	}
}
