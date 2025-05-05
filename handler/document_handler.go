package handler

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"

	"github.com/gin-gonic/gin"
	"github.com/ryangladden/archivelens-go/request"
	"github.com/ryangladden/archivelens-go/service"
)

type DocumentHandler struct {
	documentService *service.DocumentService
}

func NewDocumentHandler(documentService *service.DocumentService) *DocumentHandler {
	return &DocumentHandler{
		documentService: documentService,
	}
}

func (h *DocumentHandler) GetDocument(c *gin.Context) {

}

func (h *DocumentHandler) CreateDocument(c *gin.Context) {
	var request request.CreateDocumentRequest
	err := c.ShouldBind(&request)
	if err != nil {
		log.Error().Err(err).Msg("Error parsing form")
	}

	val := c.MustGet("user")
	userID, ok := val.(uuid.UUID)
	if !ok {
		log.Error().Msg("Unable to obtain user_id from context")
		c.AbortWithStatus(500)
		return
	}
	log.Debug().Interface("user", userID)
	request.Owner = userID
	uuid, err := h.documentService.CreateDocument(request)
	fmt.Print(uuid)
	c.JSON(200, gin.H{"form": request})
	// c.JSON(200, gin.H{"title": file.Filename})

}
