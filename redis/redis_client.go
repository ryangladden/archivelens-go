package redis

import (
	"github.com/hibiken/asynq"
	"github.com/rs/zerolog/log"
	errs "github.com/ryangladden/archivelens-go/err"
	"github.com/ryangladden/archivelens-go/microservices"
)

type RedisConnection struct {
	client *asynq.Client
}

func NewRedisConnection(endpoint string) *RedisConnection {
	return &RedisConnection{
		client: asynq.NewClient(asynq.RedisClientOpt{Addr: endpoint}),
	}
}

func (r *RedisConnection) EnqueueDocumentThumbnail(id string, filename string) error {
	task, err := microservices.NewDocumentThumbnailTask(id, filename)
	if err != nil {
		log.Error().Err(err).Msgf("Failed to enqueue document thumbnail generation for %s", id)
		return errs.ErrRedis
	}
	r.client.Enqueue(task)
	return nil
}

func (r *RedisConnection) EnqueueDocumentPreview(id string, filename string) error {
	task, err := microservices.NewDocumentPreviewTask(id, filename)
	if err != nil {
		log.Error().Err(err).Msgf("Failed to enqueue preview generation for %s", id)
		return errs.ErrRedis
	}
	r.client.Enqueue(task)
	return nil
}
