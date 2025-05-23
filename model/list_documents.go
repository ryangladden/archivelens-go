package model

import (
	"time"

	"github.com/google/uuid"
)

type ListDocumentsFilter struct {
	UserID       uuid.UUID
	Limit        int
	Page         int
	SortBy       string
	Order        string // ascending or descending
	DateMin      *time.Time
	DateMax      *time.Time
	ExcludeRoles *string
	TitleMatch   *string
	ExcludeType  *string
	Authors      *string
	IncludeTags  *string
}
