package tasks

import (
	"context"
	"encoding/json"

	"github.com/hibiken/asynq"
	"github.com/rs/zerolog/log"
)

const (
	TypeDocumentCreateThumbnail = "document:thumbnail"
	TypeDocumentCreateView      = "document:view"
)

type DocumentCreateThumbnailPayload struct {
	ID               string
	OriginalFilename string
}

type DocumentCreateViewPayload struct {
	ID               string
	OriginalFilename string
}

func NewDocumentCreateThumbnailTask(resourceID string, originalFilename string) (*asynq.Task, error) {
	payload, err := json.Marshal(DocumentCreateThumbnailPayload{
		ID:               resourceID,
		OriginalFilename: originalFilename,
	})
	if err != nil {
		return nil, err
	}
	return asynq.NewTask(TypeDocumentCreateThumbnail, payload), nil
}

func HandleDocumentCreateThumbnailTask(ctx context.Context, t *asynq.Task) error {
	var p DocumentCreateThumbnailPayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		log.Error().Err(err).Msgf("json.Unmarshal failed: %v: %s", err, asynq.SkipRetry)
	}
	log.Info().Msgf("Creating thumbnail for document %s", p.ID)
	return nil
}
