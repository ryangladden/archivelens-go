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

func (s *DocumentService) ListDocuments(request request.ListDocumentsRequest) (*response.ListDocumentsResponse, error) {
	filter := s.generateListDocumentsFilter(request)
	documentPage, err := s.documentDao.ListDocuments(filter)
	if err != nil {
		return nil, err
	}
	listResponse := s.generateListDocumentsResponse(documentPage)
	listResponse.DocumentsPerPage = filter.Limit
	listResponse.TotalDocuments = documentPage.TotalDocuments
	listResponse.PageNumber = filter.Page
	listResponse.TotalPages = int(math.Ceil(float64(documentPage.TotalDocuments) / float64(filter.Limit)))
	return listResponse, nil
}

func (s *DocumentService) generateListDocumentsResponse(page *db.DocumentPage) *response.ListDocumentsResponse {
	var listResponse response.ListDocumentsResponse
	for _, document := range page.Documents {
		inlineDocument := response.InlineDocument{
			ID:    document.Document.ID,
			Title: document.Document.Title,
			Date:  document.Document.Date,
			Author: &response.InlinePerson{
				ID:   document.DocumentMetadata.Author.ID,
				Name: document.DocumentMetadata.Author.Name,
			},
			Role: document.Document.Role,
		}
		inlineDocument.Persons, inlineDocument.Tags = s.parseSearchMetadata(document)
		listResponse.Documents = append(listResponse.Documents, inlineDocument)
	}
	return &listResponse
}

func (s *DocumentService) parseSearchMetadata(document db.InlineDocument) (*[]response.InlinePerson, *[]response.Tag) {
	var persons []response.InlinePerson
	for _, personData := range document.DocumentMetadata.Persons {
		person := response.InlinePerson{
			ID:   personData.ID,
			Name: personData.Name,
			Role: &personData.Role,
		}
		persons = append(persons, person)
	}
	var tags []response.Tag
	for _, tagData := range document.DocumentMetadata.Tags {
		tag := response.Tag{
			ID:  tagData.ID,
			Tag: tagData.Tag,
		}
		tags = append(tags, tag)
	}
	return &persons, &tags
}
