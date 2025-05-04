package model

import "github.com/google/uuid"

type Person struct {
	ID       uuid.UUID `json:"id"`
	Name     string    `json:"name"`
	Metadata string    `json:"metadata"`
}
