package service

import (
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/ryangladden/archivelens-go/db"
	errs "github.com/ryangladden/archivelens-go/err"
	"github.com/ryangladden/archivelens-go/model"
	"github.com/ryangladden/archivelens-go/request"
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
	if personModel.S3Key != "" {
		if err = s.storageManager.UploadFile(request.Avatar, personModel.S3Key); err != nil {
			return uuid.Nil, errs.ErrStorage
		}
	}
	return personModel.ID, nil
}

func (s *PersonService) generatePersonModel(request *request.CreatePersonRequest) (*model.Person, error) {
	var person model.Person
	person.Name = request.Name
	person.Birth = request.Birth
	person.Death = request.Death
	person.Summary = request.Summary

	id, err := uuid.NewV7()
	if err != nil {
		log.Error().Err(err).Msgf("Error generating uuid for new person %s", request.Name)
		return nil, err
	}

	person.ID = id
	if request.Avatar != nil {
		person.S3Key = s.storageManager.GenerateObjectKey(request.Avatar.Filename, id, "persons")
	} else {
		person.S3Key = ""
	}

	return &person, nil
}
