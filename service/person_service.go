package service

import (
	"math"
	"strings"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/ryangladden/archivelens-go/db"
	errs "github.com/ryangladden/archivelens-go/err"
	"github.com/ryangladden/archivelens-go/model"
	"github.com/ryangladden/archivelens-go/request"
	"github.com/ryangladden/archivelens-go/response"
	"github.com/ryangladden/archivelens-go/storage"
)

type PersonService struct {
	personDao      *db.PersonDAO
	storageManager *storage.StorageManager
}

func NewPersonService(personDao *db.PersonDAO, storageManager *storage.StorageManager) *PersonService {
	return &PersonService{
		personDao:      personDao,
		storageManager: storageManager,
	}
}

func (s *PersonService) CreatePerson(request *request.CreatePersonRequest) (uuid.UUID, error) {
	personModel, err := s.generatePersonModel(request)
	if err != nil {
		return uuid.Nil, errs.ErrInternalServer
	}

	if err = s.personDao.CreatePerson(personModel, request.Owner); err != nil {
		return uuid.Nil, errs.ErrDB
	}
	if personModel.S3Key != nil {
		if err = s.storageManager.UploadFile(request.Avatar, *personModel.S3Key); err != nil {
			return uuid.Nil, errs.ErrStorage
		}
	}
	return personModel.ID, nil
}

func (s *PersonService) ListPersons(request request.ListPersonsRequest) (*response.ListPersonsResponse, error) {

	filter := generateListPersonsFilter(request)
	personPage, err := s.personDao.ListPersons(filter)
	if err != nil {
		return nil, err
	}
	personList := response.ListPersonsResponse{
		PersonsPerPage: filter.Limit,
		PageNumber:     filter.Page + 1,
		TotalPersons:   personPage.TotalPersons,
		TotalPages:     int(math.Ceil(float64(personPage.TotalPersons) / float64(filter.Limit))),
	}
	for _, person := range personPage.Persons {
		personList.Persons = append(personList.Persons, *s.generatePersonResponse(person))
	}
	return &personList, nil
}

func (s *PersonService) GetPerson(request request.GetPersonRequest) (*response.PersonResponse, error) {
	person, err := s.personDao.GetPerson(request.UserID, request.PersonID)
	if err != nil {
		return nil, err
	}
	response := s.generatePersonResponse(*person)
	return response, nil
}

func (s *PersonService) generatePersonModel(request *request.CreatePersonRequest) (*model.Person, error) {
	person := model.Person{
		FirstName: &request.FirstName,
		LastName:  &request.LastName,
		Birth:     request.Birth,
		Death:     request.Death,
		Summary:   request.Summary,
	}
	id, err := uuid.NewV7()
	if err != nil {
		log.Error().Err(err).Msgf("Error generating uuid for new person %s %s", request.FirstName, request.LastName)
		return nil, err
	}

	person.ID = id
	if request.Avatar != nil {
		key := storage.GenerateObjectKey("persons", id, "avatar", request.Avatar.Filename)
		// key := fmt.Sprintf("persons/%s/original%s", id.String(), strings.ToLower(filepath.Ext(request.Avatar.Filename)))
		person.S3Key = key
	} else {
		person.S3Key = nil
	}

	return &person, nil
}

func (s *PersonService) generatePersonResponse(person model.Person) *response.PersonResponse {
	response := response.PersonResponse{
		ID:           person.ID,
		FirstName:    *person.FirstName,
		LastName:     *person.LastName,
		Birth:        person.Birth,
		Death:        person.Death,
		Summary:      person.Summary,
		Role:         *person.Role,
		PresignedUrl: s.storageManager.GeneratePresignedURL(person.S3Key),
	}
	return &response
}

func generateListPersonsFilter(request request.ListPersonsRequest) *model.ListPersonsFilter {
	nameMatch := ""
	if request.NameMatch != nil {
		nameMatch = strings.ToLower(*request.NameMatch)
	}
	filter := model.ListPersonsFilter{
		UserID:       request.UserID,
		NameMatch:    &nameMatch,
		BirthMin:     request.BirthMin,
		BirthMax:     request.BirthMax,
		DeathMin:     request.DeathMin,
		DeathMax:     request.DeathMax,
		ExcludeRoles: parseExcludeRoles(request.ExcludeRoles),
		SortBy:       parseSortBy(request.SortBy, []string{"first_name", "last_name", "birth", "death"}, "last_name"),
		Order:        parseOrder(request.Order),
	}
	filter.UserID = request.UserID
	log.Debug().Msgf("Filtering persons related to user with %s", filter.UserID)
	if request.Limit == nil {
		filter.Limit = 20
	} else {
		filter.Limit = *request.Limit
	}
	if request.Page == nil {
		filter.Page = 0
	} else {
		filter.Page = *request.Page - 1
	}
	return &filter
}
