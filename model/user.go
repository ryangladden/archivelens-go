package model

import (
	_ "github.com/google/uuid"
)

type User struct {
	ID       string `json:"id" validate:"required,uuid"`
	Email    string `json:"email" validate:"required,email"`
	Name     string `json:"name" validate:"required"`
	Password []byte `json:"password" validate:"required"`
}
