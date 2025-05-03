package db

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type ConnectionManager struct {
	DB *sql.DB
}

func NewConnectionManager(host string, port int, user string, password string, dbname string) *ConnectionManager {
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	DB, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}
	if err = DB.Ping(); err != nil {
		panic(err)
	}
	return &ConnectionManager{
		DB: DB,
	}
}
