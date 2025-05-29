package request

import (
	"mime/multipart"
	"time"

	"github.com/google/uuid"
)

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type CreateUserRequest struct {
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required,min=8"`
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
}

type CreateDocumentRequest struct {
	Title     string                `form:"title" binding:"required"`
	Type      string                `form:"type" binding:"required"`
	Author    *string               `form:"author"`
	Coauthors *[]string             `form:"coauthors"`
	Mentions  *[]string             `form:"mentions"`
	Recipient *string               `form:"recipient"`
	Date      *time.Time            `form:"date"`
	Location  *string               `form:"location" `
	File      *multipart.FileHeader `form:"file" binding:"required"`
	Owner     uuid.UUID
}

type CreatePersonRequest struct {
	FirstName string                `form:"first_name" binding:"required"`
	LastName  string                `form:"last_name" binding:"required"`
	Birth     *time.Time            `form:"birth" time_format:"2006-01-02" time_utc:"1"`
	Death     *time.Time            `form:"death" time_format:"2006-01-02" time_utc:"1"`
	Summary   *string               `form:"summary"`
	Avatar    *multipart.FileHeader `form:"file"`
	Owner     uuid.UUID
}

type ListPersonsRequest struct {
	UserID       uuid.UUID
	Page         *int       `form:"page"`
	Limit        *int       `form:"person_per_page"`
	SortBy       *string    `form:"sort_by"`
	BirthMax     *time.Time `form:"birth_max" time_format:"2006-01-02" time_utc:"1"`
	BirthMin     *time.Time `form:"birth_min" time_format:"2006-01-02" time_utc:"1"`
	DeathMax     *time.Time `form:"death_max" time_format:"2006-01-02" time_utc:"1"`
	DeathMin     *time.Time `form:"death_min" time_format:"2006-01-02" time_utc:"1"`
	NameMatch    *string    `form:"name_match"`
	ExcludeRoles *[]string  `form:"exclude_roles"`
	Order        *string    `form:"order"` // ascending or descending
}

type GetPersonRequest struct {
	UserID   uuid.UUID
	PersonID uuid.UUID
}

type ListDocumentsRequest struct {
	UserID       uuid.UUID
	Page         *int       `form:"page"`
	Limit        *int       `form:"documents_per_page"`
	SortBy       *string    `form:"sort_by"`
	DateMin      *time.Time `form:"date_min" time_format:"2006-01-02" time_utc:"1"`
	DateMax      *time.Time `form:"date_max" time_format:"2006-01-02" time_utc:"1"`
	IncludeTags  *[]string  `form:"tags"`
	TitleMatch   *string    `form:"title_match"`
	Authors      *[]string  `form:"authors"`
	ExcludeRoles *[]string  `form:"exclude_roles"`
	Order        *string    `form:"order"` // ascending or descending
	ExcludeType  *[]string  `form:"exclude_type"`
}
