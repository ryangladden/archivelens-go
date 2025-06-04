package db

// import (
// 	"context"
// 	"testing"

// 	"github.com/google/uuid"
// 	"github.com/ryangladden/archivelens-go/model"
// )

// var (
// 	cm      *ConnectionManager
// 	userDAO *UserDAO
// )

// func TestMain(m *testing.M) {
// 	cm = NewConnectionManager("localhost", 5432, "postgres", "postgres", "archive-lens-dev")
// 	userDAO = NewUserDAO(cm)
// 	m.Run()
// }

// func TestCreateUser(t *testing.T) {

// 	id, err := uuid.NewV7()
// 	if err != nil {
// 		t.Fatal("Error generating a UUID. This is likely a problem with your system")
// 	}

// 	// Create a new user
// 	user := &model.User{
// 		ID:        id,
// 		FirstName: "Test",
// 		LastName:  "User",
// 		Email:     "email@email.com",
// 		Password:  []byte("hashed-password"),
// 	}
// 	err = userDAO.CreateUser(user)
// 	if err != nil {
// 		t.Fatalf("Failed to create user: %v", err)
// 	}
// 	// Verify the user was created
// 	var createdUser model.User
// 	err = userDAO.cm.DB.QueryRow(context.Background(),
// 		"SELECT id, first_name, last_name, email, password FROM users WHERE id = $1", user.ID).Scan(
// 		&createdUser.ID, &createdUser.FirstName, &createdUser.LastName, &createdUser.Email, &createdUser.Password)
// 	if err != nil {
// 		t.Fatalf("Failed to retrieve created user: %v", err)
// 	}
// 	if createdUser.ID != user.ID || createdUser.FirstName != user.FirstName || createdUser.LastName != user.LastName || createdUser.Email != user.Email || string(createdUser.Password) != string(user.Password) {
// 		t.Errorf("Created user does not match expected user: got %+v, want %+v", createdUser, user)
// 	}
// 	// Clean up the created user
// 	_, err = userDAO.cm.DB.Exec(context.Background(), "DELETE FROM users WHERE id = $1", user.ID)
// 	if err != nil {
// 		t.Fatalf("Failed to clean up created user: %v", err)
// 	}
// }

// func TestCreateExistingEmail(t *testing.T) {
// 	id, err := uuid.NewV7()
// 	if err != nil {
// 		t.Fatalf("Failed to generate UUID: %v", err)
// 	}

// 	idNew, err := uuid.NewV7()
// 	if err != nil {
// 		t.Fatalf("Failed to generate UUID: %v", err)
// 	}

// 	// Create a new user
// 	user := &model.User{
// 		ID:        id,
// 		FirstName: "Test",
// 		LastName:  "User",
// 		Email:     "email@email.com",
// 		Password:  []byte("hashed-password"),
// 	}

// 	existingUser := &model.User{
// 		ID:        idNew,
// 		FirstName: "Existing",
// 		LastName:  "User",
// 		Email:     "email@email.com",
// 		Password:  []byte("hashed-password"),
// 	}
// 	err = userDAO.CreateUser(user)
// 	if err != nil {
// 		t.Fatalf("Failed to create user: %v", err)
// 	}
// 	err = userDAO.CreateUser(existingUser)
// 	if err == nil {
// 		t.Fatalf("Expected error when creating user with existing email, got nil")
// 	}

// 	// Verify the error message
// 	if err.Error() != "error creating user: pq: duplicate key value violates unique constraint \"users_email_key\"" {
// 		t.Errorf("Expected error message 'pq: duplicate key value violates unique constraint \"users_email_key\"', got '%v'", err)
// 	}
// 	// Clean up the created user
// 	_, err = userDAO.cm.DB.Exec(context.Background(), "DELETE FROM users WHERE id = $1", user.ID)
// 	if err != nil {
// 		t.Fatalf("Failed to clean up created user: %v", err)
// 	}
// }
