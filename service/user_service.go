package service

// import (
// 	"fmt"

// 	"github.com/go-playground/validator/v10"
// 	"github.com/rs/zerolog/log"
// 	"golang.org/x/crypto/bcrypt"

// 	"github.com/google/uuid"
// 	"github.com/ryangladden/archivelens-go/db"
// 	"github.com/ryangladden/archivelens-go/model"
// 	requests "github.com/ryangladden/archivelens-go/request"
// 	"github.com/ryangladden/archivelens-go/response"
// )

// var validate = validator.New()

// type UserService struct {
// 	userDao *db.UserDAO
// }

// func NewUserService(userDao *db.UserDAO) *UserService {
// 	return &UserService{
// 		userDao: userDao,
// 	}
// }

// func (s *UserService) CreateUser(request *requests.CreateUserRequest) (*response.LoginResponse, error) {
// 	userModel, err := createUserModel(request)
// 	log.Info().Msgf("Creating user with email: %s", request.Email)
// 	if err != nil {
// 		log.Error().Err(err).Msg("Error creating user model")
// 		return nil, fmt.Errorf("error creating user model: %w", err)
// 	}

// 	if err = s.userDao.CreateUser(userModel); err != nil {
// 		return nil, err
// 	}

// 	return &response.LoginResponse{Email: userModel.Email, FirstName: userModel.FirstName, LastName: userModel.LastName}, nil
// }

// func generateHashedPassword(password string) ([]byte, error) {
// 	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
// 	if err != nil {
// 		log.Error().Err(err).Msg("Error hashing password")
// 		return nil, fmt.Errorf("error hashing password: %w", err)
// 	}
// 	return hashedPassword, nil
// }

// func createUserModel(user *requests.CreateUserRequest) (*model.User, error) {
// 	var userModel model.User
// 	userModel.Email = user.Email
// 	userModel.FirstName = user.FirstName
// 	userModel.LastName = user.LastName
// 	hashedPassword, err := generateHashedPassword(user.Password)
// 	if err != nil {
// 		return nil, err
// 	}
// 	userModel.Password = hashedPassword

// 	id, err := uuid.NewV7()
// 	if err != nil {
// 		log.Error().Err(err).Msgf("Error generating UUID for user %s:", user.Email)
// 		return nil, fmt.Errorf("error generating UUID: %w", err)
// 	}

// 	userModel.ID = id
// 	if err := validate.Struct(userModel); err != nil {
// 		log.Error().Err(err).Msgf("Error validating user model for email %s:", user.Email)
// 		return nil, fmt.Errorf("error validating user model: %w", err)
// 	}

// 	return &userModel, nil
// }
