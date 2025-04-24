package db

import (
	_ "database/sql"
	"fmt"

	_ "github.com/lib/pq"
	"github.com/ryangladden/archivelens-go/model"
)

type UserDAO struct {
	cm *ConnectionManager
}

func NewUserDAO(cm *ConnectionManager) *UserDAO {
	return &UserDAO{cm: cm}
}

func (dao *UserDAO) CreateUser(user *model.User) error {
	_, err := dao.cm.DB.Exec(`INSERT INTO users
	(id, name, password, email)
	VALUES ($1, $2, $3, $4)`, user.ID, user.Name, user.Password, user.Email)
	if err != nil {
		return fmt.Errorf("error creating user: %v", err)
	}
	return nil
}
