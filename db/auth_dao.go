package db

import (
	"github.com/rs/zerolog/log"
	"github.com/ryangladden/archivelens-go/errs"
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
	_, err := dao.cm.DB.Exec(`INSERT INTO auth
	(token, user_id)
	VALUES ($1, $2)`, auth.AuthToken, auth.ID)
	if err != nil {
		log.Error().Err(err).Msgf("Error creating auth in database for user ID: %s", auth.ID)
		return errs.ErrDB
	}
	return nil
}
