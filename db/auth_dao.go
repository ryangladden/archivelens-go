package db

import (
	"context"
	"database/sql"

	"github.com/rs/zerolog/log"
	errs "github.com/ryangladden/archivelens-go/err"
	"github.com/ryangladden/archivelens-go/model"
)

type AuthDAO struct {
	cm *ConnectionManager
}

func NewAuthDAO(cm *ConnectionManager) *AuthDAO {
	return &AuthDAO{
		cm: cm,
	}
}

func (dao *AuthDAO) CreateAuth(auth *model.Auth) error {
	_, err := dao.cm.DB.Exec(context.Background(),
		`INSERT INTO auth
		(token, user_id)
		VALUES ($1, $2)`, auth.AuthToken, auth.ID)
	if err != nil {
		log.Error().Err(err).Msgf("Error creating auth in database for user ID: %s", auth.ID)
		return errs.ErrDB
	}
	return nil
}

func (dao *AuthDAO) GetUser(token string) (*model.User, error) {
	var user model.User

	row := dao.cm.DB.QueryRow(context.Background(),
		`SELECT user_id, first_name, last_name, email FROM auth
		INNER JOIN users ON users.id = auth.user_id
		WHERE token = $1`, token,
	)
	if err := row.Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email); err == sql.ErrNoRows {
		log.Warn().Err(err).Msgf("Failure to authenticate with token: %s", token)
		return nil, errs.ErrUnauthorized
	} else if err != nil {
		log.Error().Err(err).Msgf("Error finding auth in database for auth: %s", token)
		return nil, errs.ErrDB
	}

	return &user, nil
}

func (dao *AuthDAO) DeleteAuth(token string) error {
	_, err := dao.cm.DB.Exec(context.Background(),
		`DELETE FROM auth
		WHERE token = $1`, token,
	)
	if err == sql.ErrNoRows {
		return errs.ErrNotFound
	} else if err != nil {
		return errs.ErrDB
	}

	log.Info().Msgf("Delete auth token: %s", token)
	return nil
}
