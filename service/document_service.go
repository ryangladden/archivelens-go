package service

import (
	"fmt"
	"mime/multipart"
	"path/filepath"
	"strings"

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
	document := generateDocumentModel(request)
	s.uploadToS3(request.File, document.S3Key)
	s.documentDao.CreateDocument(request.Owner, document)
	return "", nil
}

func (s *DocumentService) uploadToS3(fileHeader *multipart.FileHeader, key string) error {
	err := s.storageManager.UploadFile(fileHeader, key)
	if err != nil {
		return err
	}
	return err
}

func generateObjectKey(originalFileName string, uuid uuid.UUID) string {
	extension := strings.ToLower(filepath.Ext(originalFileName))
	key := fmt.Sprintf(uuid.String(), extension)
	return key
}

func generateDocumentModel(request request.CreateDocumentRequest) *model.Document {
	var document model.Document
	document.Title = request.Title
	document.Location = request.Location
	document.Date = request.Date

	id, err := uuid.NewV7()
	if err != nil {
		log.Error().Err(err).Msgf("Error generating UUID for document titled \"%s\"", request.Title)
	}
	document.ID = id
	document.S3Key = generateObjectKey(request.File.Filename, id)
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
