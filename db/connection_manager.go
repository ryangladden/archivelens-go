package db

import (
	"context"
	"fmt"

	"github.com/rs/zerolog/log"

	"github.com/jackc/pgx/v5"
)

type ConnectionManager struct {
	DB *pgx.Conn
}

func NewConnectionManager(host string, port int, user string, password string, dbname string) *ConnectionManager {
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%d/%s", user, password, host, port, dbname)
	DB, err := pgx.Connect(context.Background(), connStr)
	if err != nil {
		log.Error().Err(err).Msgf("Unable to connect to database %s at host %s:%d", dbname, host, port)
		panic(err)
	}

	if err = DB.Ping(context.Background()); err != nil {
		log.Error().Err(err).Msg("Database unreachable")
		panic(err)
	}

	Init(DB)

	return &ConnectionManager{
		DB: DB,
	}
}
