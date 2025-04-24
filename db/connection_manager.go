package db

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

// var DB *sql.DB

type ConnectionManager struct {
	host     string
	port     int
	user     string
	password string
	dbname   string
	DB       *sql.DB
}

// const (
// 	host     = "localhost"
// 	port     = 5432
// 	user     = "postgres"
// 	password = "postgres"
// 	dbname   = "archive-lens-dev"
// )

func NewConnectionManager(host string, port int, user string, password string, dbname string) *ConnectionManager {
	return &ConnectionManager{
		host:     host,
		port:     port,
		user:     user,
		password: password,
		dbname:   dbname,
	}
}

func (c *ConnectionManager) Connect() error {
	var err error
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", c.host, c.port, c.user, c.password, c.dbname)
	c.DB, err = sql.Open("postgres", connStr)
	if err != nil {
		return err
	}
	return c.DB.Ping()
}
