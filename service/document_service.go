package service

import (
	"fmt"
	"math"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/ryangladden/archivelens-go/db"
	errs "github.com/ryangladden/archivelens-go/err"
	"github.com/ryangladden/archivelens-go/redis"
	"github.com/ryangladden/archivelens-go/request"
	"github.com/ryangladden/archivelens-go/response"
	"github.com/ryangladden/archivelens-go/storage"
)

type DocumentService struct {
	documentDao    *db.DocumentDAO
	storageManager *storage.StorageManager
	redisClient    *redis.RedisConnection
}

func NewDocumentService(documentDao *db.DocumentDAO, storageManager *storage.StorageManager, redisClient *redis.RedisConnection) *DocumentService {
	return &DocumentService{
		documentDao:    documentDao,
		storageManager: storageManager,
		redisClient:    redisClient,
	}
}

func (s *DocumentService) CreateDocument(request request.CreateDocumentRequest) (string, error) {
	document := s.generateDocumentModel(request)

	// Move this somewhere else
	// s3key := fmt.Sprintf("/documents/%s/original/%s", document.ID, document.OriginalFilename)
	s3key := filepath.Join("/documents", document.ID.String(), "original", document.OriginalFilename)

	err := s.storageManager.UploadMultipartFile(request.File, s3key)
	if err != nil {
		return "", errs.ErrStorage
	}

	authorships := generateAuthorshipArray(document.ID.String(), request)
	err = s.documentDao.CreateDocument(request.Owner, document, authorships)
	if err != nil {
		return "", err
	}

	err = s.redisClient.EnqueueDocumentThumbnail(document.ID.String(), document.OriginalFilename)
	if err != nil {
		return "", errs.ErrRedis
	}
	err = s.redisClient.EnqueueDocumentPreview(document.ID.String(), document.OriginalFilename)
	if err != nil {
		return "", errs.ErrRedis
	}
	err = s.redisClient.EnqueueDocumentTranscription(document.ID.String(), document.OriginalFilename)
	if err != nil {
		return "", err
	}

	return document.ID.String(), nil
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
		ID:        document.ID,
		Title:     document.Title,
		Type:      document.Type,
		Date:      document.Date,
		Location:  document.Location,
		Author:    s.generateInlinePerson(document.Author),
		Coauthors: s.generateInlinePersonList(document.Coauthors),
		Mentions:  s.generateInlinePersonList(document.Mentions),
		Recipient: s.generateInlinePerson(document.Recipient),
		Role:      document.Role,
		Tags:      document.Tags,
		Pages:     s.GetPreview(document.ID, 1, document.NumberOfPages),
	}
	return &response, nil
}

func (s *DocumentService) GetPreview(id uuid.UUID, first int, last int) []string {
	key := filepath.Join("documents", id.String(), "preview")
	var URLs []string
	for page := first; page <= last; page++ {
		pageKey := fmt.Sprintf("%s/preview-%03d.png", key, page)
		log.Debug().Msg(pageKey)
		URL := s.storageManager.GeneratePresignedURL(&pageKey)
		URLs = append(URLs, *URL)
	}
	return URLs
}
