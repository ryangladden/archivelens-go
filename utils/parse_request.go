package utils

import (
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	errs "github.com/ryangladden/archivelens-go/err"
)

func GetParamAsString(c *gin.Context, key string) string {
	val, exists := c.Params.Get(key)
	if exists {
		return val
	}
	return ""
}

func GetParamAsInt(c *gin.Context, key string, fallback int) int {
	val, exists := c.Params.Get(key)
	if exists {
		if param, err := strconv.Atoi(val); err == nil {
			return param
		}
	}
	return fallback
}

func GetParamsAsUUID(c *gin.Context, key string) (uuid.UUID, error) {
	log.Debug().Msgf("Looking for UUID with key %s", key)
	val, exists := c.Params.Get(key)
	if exists {
		if id, err := uuid.Parse(val); err == nil {
			return id, nil
		}
	}
	return uuid.Nil, errs.ErrBadRequest
}

func GetParamAsDate(c *gin.Context, key string) *time.Time {
	val, exists := c.Params.Get(key)
	if exists {
		if date, err := time.Parse(time.RFC3339, val); err == nil {
			return &date
		}
	}
	return nil
}

func GetParamsAsArray(c *gin.Context, key string) []string {
	val, exists := c.Params.Get(key)
	if exists {
		return strings.Split(val, ",")
	}
	var empty []string
	return empty
}

func GetUserIDFromContext(c *gin.Context) uuid.UUID {
	user, exists := c.Get("user")
	if exists {
		if u, ok := user.(uuid.UUID); ok {
			return u
		}
	}
	return uuid.Nil
}
