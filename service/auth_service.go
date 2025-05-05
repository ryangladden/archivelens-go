package service

import (
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/bcrypt"

	"github.com/ryangladden/archivelens-go/db"
	errs "github.com/ryangladden/archivelens-go/err"
	"github.com/ryangladden/archivelens-go/model"
	"github.com/ryangladden/archivelens-go/request"
	"github.com/ryangladden/archivelens-go/response"
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

func (s *AuthService) CreateAuth(request request.LoginRequest) (string, *response.LoginResponse, error) {
	user, err := s.userDao.GetUserByField("email", request.Email)
	if user == nil {
		log.Error().Msgf("User not found with email: %s", request.Email)
		return "", nil, errs.ErrNotFound
	}

	if err != nil {
		return "", nil, err
	}

	authModel, err := createAuthModel(user)
	if err != nil {
		log.Error().Err(err).Msg("Error creating auth model")
		return "", nil, err
	}
	err = s.verifyPassword(user, request.Password)
	if err != nil {
		log.Error().Err(err).Msg("Password verification failed")
		return "", nil, errs.ErrUnauthorized
	}

	err = s.authDao.CreateAuth(authModel)
	if err != nil {
		return "", nil, err
	}
	return authModel.AuthToken, &response.LoginResponse{Email: user.Email, Name: user.Name}, nil
}

func (s *AuthService) ValidateToken(token string) (*model.User, error) {
	user, err := s.authDao.GetUser(token)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *AuthService) DeleteAuth(token string) error {
	return s.authDao.DeleteAuth(token)
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
