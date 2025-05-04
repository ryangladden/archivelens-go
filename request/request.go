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
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
	Name     string `json:"name" binding:"required"`
}

type CreateDocumentRequest struct {
	Title     string                `form:"title" binding:"required"`
	Author    string                `form:"author"`
	Coauthors []string              `form:"coauthors"`
	Mentions  []string              `form:"mentions"`
	Recipient string                `form:"recipient"`
	Date      time.Time             `form:"date"`
	Location  string                `form:"location" `
	File      *multipart.FileHeader `form:"file" binding:"required"`
	Owner     uuid.UUID
}
