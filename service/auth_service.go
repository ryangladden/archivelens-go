package service

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/bcrypt"

	"github.com/ryangladden/archivelens-go/db"
	errs "github.com/ryangladden/archivelens-go/err"
	"github.com/ryangladden/archivelens-go/model"
	"github.com/ryangladden/archivelens-go/request"
	"github.com/ryangladden/archivelens-go/response"
)

var validate = validator.New()

type AuthService struct {
	authDao *db.AuthDAO
	// userDao *db.UserDAO
}

func NewAuthService(authDao *db.AuthDAO) *AuthService {
	return &AuthService{
		authDao: authDao,
		// userDao: userDao,
	}
}

func (s *AuthService) CreateUser(createUserRequest *request.CreateUserRequest) (string, *response.LoginResponse, error) {
	userModel, err := createUserModel(createUserRequest)
	log.Info().Msgf("Creating user with email: %s", createUserRequest.Email)
	if err != nil {
		log.Error().Err(err).Msg("Error creating user model")
		return "", nil, fmt.Errorf("error creating user model: %w", err)
	}

	if err = s.authDao.CreateUser(userModel); err != nil {
		return "", nil, err
	}

	login := request.LoginRequest{Email: userModel.Email, Password: createUserRequest.Password}
	return s.CreateAuth(login)
}

func (s *AuthService) CreateAuth(request request.LoginRequest) (string, *response.LoginResponse, error) {
	user, err := s.authDao.GetUserByField("email", request.Email)
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
	return authModel.AuthToken, &response.LoginResponse{Email: user.Email, FirstName: user.FirstName, LastName: user.LastName}, nil
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

func generateHashedPassword(password string) ([]byte, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Error().Err(err).Msg("Error hashing password")
		return nil, fmt.Errorf("error hashing password: %w", err)
	}
	return hashedPassword, nil
}

func createUserModel(user *request.CreateUserRequest) (*model.User, error) {
	var userModel model.User
	userModel.Email = user.Email
	userModel.FirstName = user.FirstName
	userModel.LastName = user.LastName
	hashedPassword, err := generateHashedPassword(user.Password)
	if err != nil {
		return nil, err
	}
	userModel.Password = hashedPassword

	id, err := uuid.NewV7()
	if err != nil {
		log.Error().Err(err).Msgf("Error generating UUID for user %s:", user.Email)
		return nil, fmt.Errorf("error generating UUID: %w", err)
	}

	userModel.ID = id
	if err := validate.Struct(userModel); err != nil {
		log.Error().Err(err).Msgf("Error validating user model for email %s:", user.Email)
		return nil, fmt.Errorf("error validating user model: %w", err)
	}

	return &userModel, nil
}
