package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	errs "github.com/ryangladden/archivelens-go/err"
	"github.com/ryangladden/archivelens-go/service"
)

func AuthenticateMiddleware(authService *service.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Debug().Msg("Middleware hit")
		token, err := c.Cookie("archive_lens_access_token")
		if err != nil {
			c.JSON(401, gin.H{"error": "Unauthorized"})
			return
		}

		user, err := authService.ValidateToken(token)
		if err == errs.ErrDB {
			c.AbortWithStatus(500)
			return
		} else if err == errs.ErrUnauthorized {
			c.AbortWithStatus(401)
			return
		}
		log.Info().Msgf("User: %s", user.Name)
		log.Debug().Msgf("Middleware: setting user_id %s", user.ID)
		c.Set("user_id", user.ID)
		c.Next()
	}
}
