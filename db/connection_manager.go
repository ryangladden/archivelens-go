package db

import (
	"database/sql"
	"fmt"

	"github.com/rs/zerolog/log"

	_ "github.com/lib/pq"
)

type ConnectionManager struct {
	DB *sql.DB
}

func NewConnectionManager(host string, port int, user string, password string, dbname string) *ConnectionManager {
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	DB, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Error().Err(err).Msgf("Unable to connect to database %s at host %s:%d", dbname, host, port)
		panic(err)
	}

	if err = DB.Ping(); err != nil {
		log.Error().Err(err).Msg("Database unreachable")
		panic(err)
	}

	Init(DB)

	return &ConnectionManager{
		DB: DB,
	}
}
