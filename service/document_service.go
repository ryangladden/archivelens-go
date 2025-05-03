package service

import (
	"mime/multipart"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/ryangladden/archivelens-go/db"
	"github.com/ryangladden/archivelens-go/model"
	"github.com/ryangladden/archivelens-go/request"
	"github.com/ryangladden/archivelens-go/storage"
)

type DocumentService struct {
	documentDao *db.DocumentDAO
	storage     *storage.StorageManager
}

func NewDocumentService(documentDao *db.DocumentDAO) *DocumentService {
	return &DocumentService{
		documentDao: documentDao,
	}
}

func (s *DocumentService) CreateDocument(request request.CreateDocumentRequest) (string, error) {
	var document model.Document
	document.Title = request.Title
	document.Location = request.Location
	document.Date = request.Date
	id, err := uuid.NewV7()
	if err != nil {
		log.Error().Err(err).Msgf("Error generating UUID for document titled \"%s\"", request.Title)
	}
	document.ID = id.String()
	s.uploadToS3(request.File, id.String())
	return "", nil
}

func (s *DocumentService) uploadToS3(fileHeader *multipart.FileHeader, id string) error {
	err := s.storage.UploadFile(fileHeader, "archive-lens", "sample.pdf")
	if err != nil {
		return err
	}
	return err
}
