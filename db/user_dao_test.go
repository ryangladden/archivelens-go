package db

import (
	"testing"

	"github.com/google/uuid"
	"github.com/ryangladden/archivelens-go/model"
)

var (
	cm      *ConnectionManager
	userDAO *UserDAO
)

func TestMain(m *testing.M) {
	cm = NewConnectionManager("localhost", 5432, "postgres", "postgres", "archive-lens-dev")
	userDAO = NewUserDAO(cm)
	m.Run()
}

func TestCreateUser(t *testing.T) {

	id, err := uuid.NewV7()

	// Create a new user
	user := &model.User{
		ID:       id,
		Name:     "Test User",
		Email:    "email@email.com",
		Password: []byte("hashed-password"),
	}
	err = userDAO.CreateUser(user)
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}
	// Verify the user was created
	var createdUser model.User
	err = userDAO.cm.DB.QueryRow(
		"SELECT id, name, email, password FROM users WHERE id = $1", user.ID).Scan(
		&createdUser.ID, &createdUser.Name, &createdUser.Email, &createdUser.Password)
	if err != nil {
		t.Fatalf("Failed to retrieve created user: %v", err)
	}
	if createdUser.ID != user.ID || createdUser.Name != user.Name || createdUser.Email != user.Email || string(createdUser.Password) != string(user.Password) {
		t.Errorf("Created user does not match expected user: got %+v, want %+v", createdUser, user)
	}
	// Clean up the created user
	_, err = userDAO.cm.DB.Exec("DELETE FROM users WHERE id = $1", user.ID)
	if err != nil {
		t.Fatalf("Failed to clean up created user: %v", err)
	}
}

func TestCreateExistingEmail(t *testing.T) {
	id, err := uuid.NewV7()
	if err != nil {
		t.Fatalf("Failed to generate UUID: %v", err)
	}

	idNew, err := uuid.NewV7()
	if err != nil {
		t.Fatalf("Failed to generate UUID: %v", err)
	}

	// Create a new user
	user := &model.User{
		ID:       id,
		Name:     "Test User",
		Email:    "email@email.com",
		Password: []byte("hashed-password"),
	}

	existingUser := &model.User{
		ID:       idNew,
		Name:     "Existing User",
		Email:    "email@email.com",
		Password: []byte("hashed-password"),
	}
	err = userDAO.CreateUser(user)
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}
	err = userDAO.CreateUser(existingUser)
	if err == nil {
		t.Fatalf("Expected error when creating user with existing email, got nil")
	}

	// Verify the error message
	if err.Error() != "error creating user: pq: duplicate key value violates unique constraint \"users_email_key\"" {
		t.Errorf("Expected error message 'pq: duplicate key value violates unique constraint \"users_email_key\"', got '%v'", err)
	}
	// Clean up the created user
	_, err = userDAO.cm.DB.Exec("DELETE FROM users WHERE id = $1", user.ID)
	if err != nil {
		t.Fatalf("Failed to clean up created user: %v", err)
	}
}
