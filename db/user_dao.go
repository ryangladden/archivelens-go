package db

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
	"github.com/ryangladden/archivelens-go/model"
)

type UserDAO struct {
	db *sql.DB
}

func NewUserDAO(db *sql.DB) *UserDAO {
	return &UserDAO{db: db}
}

func (h *UserDAO) CreateUser(user *model.User) error {
	_, err := DB.Exec(`INSERT INTO users
	(id, name, password, email)
	VALUES ($1, $2, $3, $4)`, user.ID, user.Name, user.Password, user.Email)
	if err != nil {
		return fmt.Errorf("error creating user: %v", err)
	}
	return nil
}
