package response

import (
	"time"

	"github.com/google/uuid"
)

type LoginResponse struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
}

type CreatePersonResonse struct {
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	ID        uuid.UUID `json:"id"`
}

type PersonResponse struct {
	ID           uuid.UUID  `json:"id"`
	FirstName    string     `json:"first_name"`
	LastName     string     `json:"last_name"`
	Birth        *time.Time `json:"birth"`
	Death        *time.Time `json:"death" omitempty`
	Summary      *string    `json:"summary"`
	PresignedUrl *string    `json:"avatar"`
	Role         string     `json:"role"`
}

type ListPersonsResponse struct {
	Persons        []PersonResponse `json:"persons"`
	PageNumber     int              `json:"page"`
	TotalPages     int              `json:"total_pages"`
	PersonsPerPage int              `json:"persons_per_page"`
	TotalPersons   int              `json:"total_persons"`
}
