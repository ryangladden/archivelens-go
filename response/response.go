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
	Birth        *time.Time `json:"birth" time_format:"2006-01-02" time_utc:"1"`
	Death        *time.Time `json:"death" time_format:"2006-01-02" time_utc:"1"`
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

type InlinePerson struct {
	ID           uuid.UUID `json:"id"`
	Name         string    `json:"name"`
	PresignedURL *string   `json:"avatar"`
}

type DocumentResponse struct {
	ID           uuid.UUID       `json:"id"`
	Title        string          `json:"title"`
	Type         string          `json:"type"`
	Date         *time.Time      `json:"date"`
	Location     *string         `json:"location"`
	Author       *InlinePerson   `json:"author"`
	Coauthors    *[]InlinePerson `json:"coauthors"`
	Mentions     *[]InlinePerson `json:"mentions"`
	Recipient    *InlinePerson   `json:"recipient"`
	Role         string          `json:"role"`
	PresignedUrl string          `json:"src"`
}

type InlineDocument struct {
	ID        uuid.UUID     `json:"id"`
	Title     string        `json:"title"`
	Date      *time.Time    `json:"date"`
	Type      string        `json:"type"`
	Author    *InlinePerson `json:"author"`
	Thumbnail string        `json:"thumbnail"`
	Role      string        `json:"role"`
}

type ListDocumentsResponse struct {
	Documents        []InlineDocument `json:"documents"`
	PageNumber       int              `json:"page"`
	TotalPages       int              `json:"total_pages"`
	DocumentsPerPage int              `json:"persons_per_page"`
	TotalDocuments   int              `json:"total_persons"`
}
