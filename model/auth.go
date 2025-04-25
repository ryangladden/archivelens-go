package model

type Auth struct {
	ID        string `json:"id" validate:"required,uuid"`
	AuthToken string `json:"auth_token" validate:"required,uuid"`
}
