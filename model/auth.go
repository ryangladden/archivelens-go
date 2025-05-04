package model

import "github.com/google/uuid"

type Auth struct {
	ID        uuid.UUID `json:"id" validate:"required,uuid"`
	AuthToken string    `json:"auth_token" validate:"required,uuid"`
}
