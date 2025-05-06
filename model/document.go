package model

import (
	"time"

	"github.com/google/uuid"
)

type Document struct {
	ID       uuid.UUID  `json:"id`
	Title    string     `json:"title"`
	Date     *time.Time `json:"date"`
	Location *string    `json:"location"`
	S3Key    string     `json:"s3key"`
}
