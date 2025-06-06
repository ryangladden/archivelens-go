package model

import (
	"time"

	"github.com/google/uuid"
)

type ListPersonsFilter struct {
	UserID       uuid.UUID
	Limit        int
	Page         int
	SortBy       string
	Order        string // ascending or descending
	BirthMin     *time.Time
	BirthMax     *time.Time
	DeathMin     *time.Time
	DeathMax     *time.Time
	ExcludeRoles *string
	NameMatch    *string
}
