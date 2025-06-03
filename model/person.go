package model

import (
	"time"

	"github.com/google/uuid"
)

type Person struct {
	ID        uuid.UUID  `json:"id"`
	FirstName *string    `json:"first_name" validate:"required"`
	LastName  *string    `json:"last_name" validate:"required"`
	S3Key     *string    `json:"s3key"`
	Birth     *time.Time `json:"birth"`
	Death     *time.Time `json:"death"`
	Summary   *string    `json:"summary"`
	Role      *string
}
