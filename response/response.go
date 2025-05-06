package response

import "github.com/google/uuid"

type LoginResponse struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type CreatePersonResonse struct {
	Name string    `json:"name"`
	ID   uuid.UUID `json:"id"`
}
