package model

import (
	"time"
)

type Document struct {
	ID       string    `json:"id validate:"uuid"`
	Title    string    `json:"title"`
	Date     time.Time `json:"date"`
	Location string    `json:"location"`
	S3Key    string    `json:"s3key"`
}
