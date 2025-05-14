package db

import (
	"context"
	"fmt"

	"github.com/lib/pq"
	"github.com/rs/zerolog/log"
	errs "github.com/ryangladden/archivelens-go/err"
	"github.com/ryangladden/archivelens-go/model"
)

type UserDAO struct {
	cm *ConnectionManager
}

func NewUserDAO(cm *ConnectionManager) *UserDAO {
	return &UserDAO{cm: cm}
}

func (dao *UserDAO) CreateUser(user *model.User) error {
	_, err := dao.cm.DB.Exec(context.Background(), `INSERT INTO users
	(id, first_name, last_name, password, email)
	VALUES ($1, $2, $3, $4, $5)`, user.ID, user.FirstName, user.LastName, user.Password, user.Email)
	if err, ok := err.(*pq.Error); ok && err.Code == "23505" {
		log.Error().Err(err).Msg("User with this email already exists")
		return errs.ErrConflict
	} else if err != nil {
		log.Error().Err(err).Msg("Error creating user in database")
		return errs.ErrDB
		// return fmt.Errorf("error creating user in database: %v", err)
	}
	return nil
}

func (dao *UserDAO) GetUserByField(field string, value string) (*model.User, error) {
	var user model.User
	query := fmt.Sprintf("SELECT id, first_name, last_name, email, password FROM users WHERE %s = $1", field)
	err := dao.cm.DB.QueryRow(context.Background(), query, value).Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.Password)
	if err != nil {
		log.Error().Err(err).Msgf("error retrieving user by %s", field)
		return nil, fmt.Errorf("error retrieving user by %s: %v", field, err)
	}
	return &user, nil
}
