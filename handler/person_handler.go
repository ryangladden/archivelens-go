package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/ryangladden/archivelens-go/request"
	"github.com/ryangladden/archivelens-go/response"
	"github.com/ryangladden/archivelens-go/service"
)

type PersonHandler struct {
	personService *service.PersonService
}

func NewPersonHandler(personService *service.PersonService) *PersonHandler {
	return &PersonHandler{
		personService: personService,
	}
}

func (h *PersonHandler) CreatePerson(c *gin.Context) {
	var request request.CreatePersonRequest

	if err := c.ShouldBind(&request); err != nil {
		log.Error().Err(err).Msg("Invalid create person request")
		c.AbortWithStatusJSON(400, gin.H{"error": "invalid request body"})
		return
	}

	file, err := c.FormFile("file")
	if err != nil {
		if err == http.ErrMissingFile {
			request.Avatar = nil
		} else {
			log.Error().Err(err).Msg("Error getting avatar from form")
			request.Avatar = nil
		}
	} else {
		request.Avatar = file
	}
	val := c.MustGet("user")
	userID, ok := val.(uuid.UUID)
	if !ok {
		log.Error().Msg("Unable to obtain user_id from context")
		c.AbortWithStatus(500)
		return
	}
	request.Owner = userID

	id, err := h.personService.CreatePerson(&request)
	if err != nil {
		c.AbortWithStatus(500)
		return
	}

	response := response.CreatePersonResonse{Name: request.Name, ID: id}
	c.JSON(201, response)
}
