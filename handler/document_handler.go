package handler

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"

	"github.com/gin-gonic/gin"
	"github.com/ryangladden/archivelens-go/request"
	"github.com/ryangladden/archivelens-go/service"
	"github.com/ryangladden/archivelens-go/utils"
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

	log.Debug().Msg("Get document called")
	val := c.MustGet("user")
	userID, ok := val.(uuid.UUID)
	if !ok {
		log.Error().Msg("Unable to obtain user_id from context")
		c.AbortWithStatus(403)
		return
	}

	request := request.GetDocumentRequest{
		UserID: userID,
	}
	var err error
	request.DocumentID, err = utils.GetParamsAsUUID(c, "id")
	if err != nil {
		log.Error().Err(err).Msgf("Failed to get document by id")
		c.AbortWithStatus(400)
		return
	}
	document, err := h.documentService.GetDocument(request)
	if err != nil {
		c.AbortWithStatus(400)
		return
	}
	fmt.Println(document)
	c.JSON(200, document)
}

func (h *DocumentHandler) CreateDocument(c *gin.Context) {
	var request request.CreateDocumentRequest
	err := c.ShouldBind(&request)
	if err != nil {
		log.Error().Err(err).Msg("Error parsing form")
		c.AbortWithStatus(400)
		return
	}

	if !utils.ValidateMIMEType(request.File, []string{"application/pdf", "image/png", "image/jpeg"}) {
		c.AbortWithStatus(400)
	}

	val := c.MustGet("user")
	userID, ok := val.(uuid.UUID)
	if !ok {
		log.Error().Msg("Unable to obtain user_id from context")
		c.AbortWithStatus(403)
		return
	}
	log.Debug().Interface("user", userID)
	request.Owner = userID
	uuid, err := h.documentService.CreateDocument(request)
	if err != nil {
		c.AbortWithStatus(500)
		return
	}
	c.JSON(200, gin.H{"id": uuid})
	// c.JSON(200, gin.H{"title": file.Filename})

}

func (h *DocumentHandler) ListDocuments(c *gin.Context) {
	var request request.ListDocumentsRequest
	err := c.ShouldBind(&request)
	if err != nil {
		log.Error().Err(err).Msg("Invalid query for listing documents")
		c.AbortWithStatus(400)
		return
	}

	request.UserID = utils.GetUserIDFromContext(c)
	documents, err := h.documentService.ListDocuments(request)
	if err != nil {
		c.AbortWithStatus(500)
		return
	}
	c.JSON(200, documents)
}

func (h *DocumentHandler) GetPreview(c *gin.Context) {
	log.Debug().Msg("Get preview called")
	id, err := utils.GetParamsAsUUID(c, "id")
	if err != nil {
		log.Error().Err(err).Msg("Invalid UUID")
		c.AbortWithStatus(400)
		return
	}
	first := utils.GetParamAsInt(c, "first", 1)
	last := utils.GetParamAsInt(c, "last", 10)
	log.Debug().Msgf("Preview requested for %s, first page: %d, last page: %d", id, first, last)

	previewURLs := h.documentService.GetPreview(id, first, last)
	c.JSON(200, previewURLs)
}
