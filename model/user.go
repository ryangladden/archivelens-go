package model

import (
	"github.com/google/uuid"
)

type User struct {
	ID       uuid.UUID `json:"id" validate:"required,uuid"`
	Email    string    `json:"email" validate:"required,email"`
	Name     string    `json:"name" validate:"required"`
	Password []byte    `json:"password" validate:"required"`
}
