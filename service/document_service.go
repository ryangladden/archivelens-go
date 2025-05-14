package service

import (
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/ryangladden/archivelens-go/db"
	"github.com/ryangladden/archivelens-go/model"
	"github.com/ryangladden/archivelens-go/request"
	"github.com/ryangladden/archivelens-go/storage"
)

type DocumentService struct {
	documentDao    *db.DocumentDAO
	storageManager *storage.StorageManager
}

func NewDocumentService(documentDao *db.DocumentDAO, storageManager *storage.StorageManager) *DocumentService {
	return &DocumentService{
		documentDao:    documentDao,
		storageManager: storageManager,
	}
}

func (s *DocumentService) CreateDocument(request request.CreateDocumentRequest) (string, error) {
	document := s.generateDocumentModel(request)
	err := s.storageManager.UploadFile(request.File, &document.S3Key)
	if err != nil {
		return "", err
	}
	authorships := generateAuthorshipArray(document.ID.String(), request)
	err = s.documentDao.CreateDocument(request.Owner, document, authorships)
	if err != nil {
		return "", err
	}
	return "", nil
}

func (s *DocumentService) generateDocumentModel(request request.CreateDocumentRequest) *model.Document {
	var document model.Document
	document.Title = request.Title
	document.Location = request.Location
	document.Date = request.Date

	id, err := uuid.NewV7()
	if err != nil {
		log.Error().Err(err).Msgf("Error generating UUID for document titled \"%s\"", request.Title)
	}
	document.ID = id
	s3Key := s.storageManager.GenerateObjectKey(request.File.Filename, id, "documents")
	document.S3Key = *s3Key
	return &document
}

func createAuthorship(personIds []string, documentId string, role string) []model.Authorship {
	var authorships []model.Authorship
	for _, id := range personIds {
		authorships = append(authorships, model.Authorship{
			PersonID:   id,
			DocumentID: documentId,
			Role:       role,
		})
	}
	return authorships
}

func generateAuthorshipArray(documentId string, request request.CreateDocumentRequest) []model.Authorship {
	var authorships []model.Authorship
	authorships = append(authorships, createAuthorship([]string{request.Author}, documentId, "author")...)
	authorships = append(authorships, createAuthorship(request.Coauthors, documentId, "coauthor")...)
	authorships = append(authorships, createAuthorship(request.Mentions, documentId, "mentioned")...)
	authorships = append(authorships, createAuthorship([]string{request.Recipient}, documentId, "recipient")...)
	return authorships
}
