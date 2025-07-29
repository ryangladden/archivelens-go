package service

import (
	"math"

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
	err := s.storageManager.UploadFile(request.File, document.S3Key)
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

func (s *DocumentService) ListDocuments(request request.ListDocumentsRequest) (*response.ListDocumentsResponse, error) {
	filter := s.generateListDocumentsFilter(request)
	documentPage, err := s.documentDao.ListDocuments(filter)
	if err != nil {
		return nil, err
	}
	listResponse := s.generateListDocumentsResponse(documentPage)
	listResponse.DocumentsPerPage = filter.Limit
	listResponse.TotalDocuments = documentPage.TotalDocuments
	listResponse.PageNumber = filter.Page + 1
	listResponse.TotalPages = int(math.Ceil(float64(documentPage.TotalDocuments) / float64(filter.Limit)))
	return listResponse, nil
}

func (s *DocumentService) GetDocument(request request.GetDocumentRequest) (*response.DocumentResponse, error) {

	document, err := s.documentDao.GetDocument(request.UserID, request.DocumentID)
	if err != nil {
		return nil, err
	}
	response := response.DocumentResponse{
		ID:           document.ID,
		Title:        document.Title,
		Type:         document.Type,
		Date:         document.Date,
		Location:     document.Location,
		Author:       s.generateInlinePerson(document.Author),
		Coauthors:    s.generateInlinePersonList(document.Coauthors),
		Mentions:     s.generateInlinePersonList(document.Mentions),
		Recipient:    s.generateInlinePerson(document.Recipient),
		Role:         document.Role,
		PresignedUrl: *s.storageManager.GeneratePresignedURL(&document.S3Key),
		Tags:         document.Tags,
	}
	return &response, nil
}
