package service

import (
	"fmt"

	"github.com/ryangladden/archivelens-go/db"
	"github.com/ryangladden/archivelens-go/request"
	"github.com/ryangladden/archivelens-go/response"
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

func (s *DocumentService) ListDocuments(request request.ListDocumentsRequest) ([]response.ListDocumentsResponse, error) {
	filter := s.generateListDocumentsFilter(request)
	documentPage, err := s.documentDao.ListDocuments(filter)
	if err != nil {
		return nil, err
	}
	fmt.Println(documentPage.Documents[len(documentPage.Documents)-1].Title)
	return []response.ListDocumentsResponse{}, nil
}
