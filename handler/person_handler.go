package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/ryangladden/archivelens-go/request"
	"github.com/ryangladden/archivelens-go/response"
	"github.com/ryangladden/archivelens-go/service"
	"github.com/ryangladden/archivelens-go/utils"
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

	userID := utils.GetUserIDFromContext(c)
	if userID == uuid.Nil {
		log.Error().Msg("Unable to obtain user_id from context")
		c.AbortWithStatus(403)
		return
	}
	log.Debug().Msgf("User ID: %s", userID.String())
	request.Owner = userID

	id, err := h.personService.CreatePerson(&request)
	if err != nil {
		c.AbortWithStatus(500)
		return
	}

	response := response.CreatePersonResonse{FirstName: request.FirstName, LastName: request.LastName, ID: id}
	c.JSON(201, response)
}

func (h *PersonHandler) ListPersons(c *gin.Context) {

	var request request.ListPersonsRequest
	err := c.ShouldBind(&request)
	if err != nil {
		log.Error().Err(err).Msg("Invalid query for list persons")
		c.AbortWithStatus(400)
		return
	}
	request.UserID = utils.GetUserIDFromContext(c)
	persons, err := h.personService.ListPersons(request)
	if err != nil {
		c.AbortWithStatus(500)
		return
	}
	// c.JSON(200, persons)
	c.JSON(200, persons)
}

// func createListRequestFromParams(c *gin.Context) (*request.ListPersonsRequest, error) {
// 	var request request.ListPersonsRequest
// 	request.UserID = utils.GetUserIDFromContext(c)
// 	if request.UserID == uuid.Nil {
// 		log.Error().Msg("Unable to obtain user_id from context")
// 		c.AbortWithStatus(403)
// 		return nil, errs.ErrForbidden
// 	}
// 	request.Page = utils.GetParamAsInt(c, "page", 0)
// 	request.Limit = utils.GetParamAsInt(c, "limit", 20)
// 	request.SortBy = parseSortBy(c)
// 	request.BirthMax = utils.GetParamAsDate(c, "birth_max")
// 	request.BirthMin = utils.GetParamAsDate(c, "birth_min")
// 	request.DeathMax = utils.GetParamAsDate(c, "death_max")
// 	request.DeathMin = utils.GetParamAsDate(c, "death_min")
// 	request.ExcludeRoles = parseExcludeRoles(c)
// 	return &request, nil
// }

// func parseSortBy(c *gin.Context) string {
// 	val := utils.GetParamAsString(c, "sortby")
// 	switch val {
// 	case "birth", "death", "first_name":
// 		return val
// 	}
// 	return "last_name"
// }

// func parseExcludeRoles(c *gin.Context) []string {
// 	params := utils.GetParamsAsArray(c, "exclude_roles")
// 	roles := make([]string, 0, len(params))
// 	for _, role := range params {
// 		if slices.Contains([]string{"owner", "viewer", "editor"}, role) {
// 			roles = append(roles, role)
// 		}
// 	}
// 	return roles
// }
