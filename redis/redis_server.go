package redis

import (
	"github.com/hibiken/asynq"
	"github.com/ryangladden/archivelens-go/db"
	"github.com/ryangladden/archivelens-go/microservices"
	"github.com/ryangladden/archivelens-go/storage"
)

type RedisWorker struct {
	redisServer    *asynq.Server
	mux            *asynq.ServeMux
	documentWorker *microservices.DocumentWorker
}

func NewRedisWorker(endpoint string, storageManager *storage.StorageManager, documentDAO *db.DocumentDAO) *RedisWorker {
	redisServer := asynq.NewServer(
		asynq.RedisClientOpt{Addr: endpoint},
		asynq.Config{Concurrency: 10},
	)
	documentWorker := microservices.NewDocumentWorker(documentDAO, storageManager)
	mux := asynq.NewServeMux()

	redisWorker := RedisWorker{
		redisServer:    redisServer,
		mux:            mux,
		documentWorker: documentWorker,
	}

	redisWorker.addHandlers()
	go redisServer.Run(mux)
	return &redisWorker
}

func (rw *RedisWorker) addHandlers() {
	rw.mux.HandleFunc(microservices.TypeDocumentThumbnail, rw.documentWorker.HandleDocumentThumbnailTask)
	rw.mux.HandleFunc(microservices.TypeDocumentPreview, rw.documentWorker.HandleDocumentPreviewTask)
}
