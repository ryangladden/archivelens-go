package service

import (
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/bcrypt"

	"github.com/ryangladden/archivelens-go/db"
	"github.com/ryangladden/archivelens-go/errs"
	"github.com/ryangladden/archivelens-go/model"
	"github.com/ryangladden/archivelens-go/requests"
)

type AuthService struct {
	authDao *db.AuthDAO
	userDao *db.UserDAO
}

func NewAuthService(authDao *db.AuthDAO, userDao *db.UserDAO) *AuthService {
	return &AuthService{
		authDao: authDao,
		userDao: userDao,
	}
}

func (s *AuthService) CreateAuth(request requests.LoginRequest) (string, error) {
	user, err := s.userDao.GetUserByField("email", request.Email)
	if user == nil {
		log.Error().Msgf("User not found with email: %s", request.Email)
		return "", errs.ErrNotFound
	}

	if err != nil {
		return "", err
	}

	authModel, err := createAuthModel(user)
	if err != nil {
		log.Error().Err(err).Msg("Error creating auth model")
		return "", err
	}

	err = s.verifyPassword(user, request.Password)
	if err != nil {
		log.Error().Err(err).Msg("Password verification failed")
		return "", errs.ErrUnauthorized
	}

	err = s.authDao.CreateAuth(authModel)
	if err != nil {
		return "", err
	}
	return authModel.AuthToken, nil
}

func createAuthModel(user *model.User) (*model.Auth, error) {
	authModel := &model.Auth{
		ID:        user.ID,
		AuthToken: uuid.NewString(),
	}
	if err := validate.Struct(authModel); err != nil {
		return nil, err
	}
	return authModel, nil
}

func (s *AuthService) verifyPassword(user *model.User, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
}
