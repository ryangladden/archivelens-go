package model

import (
	"time"

	"github.com/google/uuid"
)

type Document struct {
	ID        uuid.UUID  `json:"id"`
	Title     string     `json:"title"`
	Date      *time.Time `json:"date"`
	Location  *string    `json:"location"`
	Type      string     `json:"type"`
	S3Key     string     `json:"s3key"`
	Author    *Person
	Coauthors *[]Person
	Mentions  *[]Person
	Recipient *Person
	Role      string
	Tags      *[]Tag
}

type Tag struct {
	ID  int
	Tag string
}
