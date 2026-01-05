package microservices

import (
	"context"
	"encoding/json"
	"fmt"
	"path/filepath"
	"slices"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/hibiken/asynq"
	"github.com/rs/zerolog/log"
	"github.com/ryangladden/archivelens-go/db"
	"github.com/ryangladden/archivelens-go/storage"
)

type DocumentWorker struct {
	redisConnection *asynq.Server
	documentDao     *db.DocumentDAO
	storageManager  *storage.StorageManager
}

func NewDocumentWorker(documentDao *db.DocumentDAO, storageManager *storage.StorageManager) *DocumentWorker {
	return &DocumentWorker{
		documentDao:    documentDao,
		storageManager: storageManager,
	}
}

const (
	TypeDocumentThumbnail         = "document:thumbnail"
	TypeDocumentPreview           = "document:preview"
	TypeDocumentTranscribeAudio   = "document:transcribe:audio"
	TypeDocumentTranscribeWritten = "document:transcribe:htr"
)

var (
	WrittenDocuments = []string{".pdf", ".jpg", ".jpeg", ".png"}
	AudioDocuments   = []string{".wav", ".mp3", ".m4a", ".aac", ".ogg", ".opus"}
)

type DocumentPayload struct {
	ID               string
	OriginalFilename string
}

type DocumentProcessor struct {
	storageManager *storage.StorageManager
}

func NewDocumentProcessor(storageManager *storage.StorageManager) *DocumentProcessor {
	return &DocumentProcessor{
		storageManager: storageManager,
	}
}

func NewDocumentThumbnailTask(resourceID string, originalFilename string) (*asynq.Task, error) {
	payload, err := marshalPayload(resourceID, originalFilename)
	if err != nil {
		return nil, err
	}
	return asynq.NewTask(TypeDocumentThumbnail, payload), nil
}

func (dw *DocumentWorker) HandleDocumentThumbnailTask(ctx context.Context, t *asynq.Task) error {

	p, err := dw.unmarshalPayload(t)
	if err != nil {
		return err
	}
	err = dw.documentDao.UpdateDocumentJobStatus(uuid.MustParse(p.ID), "thumbnail", "processing")
	if err != nil {
		return err
	}

	log.Info().Msgf("Creating thumbnail for document %s", p.ID)
	err = dw.GenerateThumb(p.ID, p.OriginalFilename)
	if err != nil {
		return err
	}
	err = dw.documentDao.UpdateDocumentJobStatus(uuid.MustParse(p.ID), "thumbnail", "processed")
	if err != nil {
		return err
	}
	return nil
}

func NewDocumentPreviewTask(resourceID string, originalFilename string) (*asynq.Task, error) {
	payload, err := marshalPayload(resourceID, originalFilename)
	if err != nil {
		return nil, err
	}
	return asynq.NewTask(TypeDocumentPreview, payload), nil
}

func (dw *DocumentWorker) HandleDocumentPreviewTask(ctx context.Context, t *asynq.Task) error {
	p, err := dw.unmarshalPayload(t)
	if err != nil {
		return err
	}

	err = dw.documentDao.UpdateDocumentJobStatus(uuid.MustParse(p.ID), "preview", "processing")
	if err != nil {
		return err
	}

	log.Info().Msgf("Generating preview for document %s", p.ID)
	pages, err := dw.GeneratePreview(p.ID, p.OriginalFilename)
	if err != nil {
		return err
	}
	log.Debug().Msgf("Document %s has %d pages", p.ID, pages)

	dw.documentDao.UpdateDocument(uuid.MustParse(p.ID), "pages", strconv.Itoa(pages))

	err = dw.documentDao.UpdateDocumentJobStatus(uuid.MustParse(p.ID), "preview", "processed")
	if err != nil {
		return err
	}
	return nil
}

func NewDocumentTranscriptionTask(resourceID string, originalFilename string) (*asynq.Task, error) {
	payload, err := marshalPayload(resourceID, originalFilename)
	if err != nil {
		return nil, err
	}

	extension := strings.ToLower(filepath.Ext(originalFilename))
	if slices.Contains(WrittenDocuments, extension) {
		return asynq.NewTask(TypeDocumentTranscribeWritten, payload), nil
	} else if slices.Contains(AudioDocuments, extension) {
		return asynq.NewTask(TypeDocumentTranscribeAudio, payload), nil
	}

	return nil, fmt.Errorf("unsupported file extension: %s", extension)
}

func marshalPayload(resourceID string, originalFilename string) ([]byte, error) {
	payload, err := json.Marshal(DocumentPayload{
		ID:               resourceID,
		OriginalFilename: originalFilename,
	})
	if err != nil {
		log.Error().Err(err).Msgf("Failed to generate document payload for resource %s", resourceID)
		return nil, err
	}

	return payload, nil
}

func (dw *DocumentWorker) unmarshalPayload(t *asynq.Task) (*DocumentPayload, error) {
	var p DocumentPayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		log.Error().Err(err).Msgf("json.Unmarshal failed: %v: %s", err, asynq.SkipRetry)
		return nil, err
	}
	return &p, nil
}
