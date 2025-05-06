package model

import (
	"time"

	"github.com/google/uuid"
)

type Person struct {
	ID      uuid.UUID  `json:"id"`
	Name    string     `json:"name" validate:"required"`
	S3Key   string     `json:"s3key`
	Birth   *time.Time `json:"birth"`
	Death   *time.Time `json:"death"`
	Summary *string    `json:"summary"`
}
