package tasks

import (
	"context"
	"encoding/json"

	"github.com/hibiken/asynq"
	"github.com/rs/zerolog/log"
)

const (
	TypeDocumentThumbnail         = "document:thumbnail"
	TypeDocumentView              = "document:view"
	TypeDocumentTranscribeAudio   = "document:transcribe:audio"
	TypeDocumentTranscribeWritten = "document:transcribe:htr"
)

type DocumentPayload struct {
	ID               string
	OriginalFilename string
}

func NewDocumentThumbnailTask(resourceID string, originalFilename string) (*asynq.Task, error) {
	payload, err := marshalPayload(resourceID, originalFilename)
	if err != nil {
		return nil, err
	}
	return asynq.NewTask(TypeDocumentThumbnail, payload), nil
}

func HandleDocumentThumbnailTask(ctx context.Context, t *asynq.Task) error {
	p, err := unmarshalPayload(t)
	if err != nil {
		return err
	}
	log.Info().Msgf("Creating thumbnail for document %s", p.ID)
	return nil
}

func NewDocumentPreviewTask(resourceID string, originalFilename string) (*asynq.Task, error) {
	return nil, nil
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

func unmarshalPayload(t *asynq.Task) (*DocumentPayload, error) {
	var p DocumentPayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		log.Error().Err(err).Msgf("json.Unmarshal failed: %v: %s", err, asynq.SkipRetry)
		return nil, err
	}
	return &p, nil
}
