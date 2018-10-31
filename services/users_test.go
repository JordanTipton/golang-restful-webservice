package services_test

import (
	"fmt"
	"testing"

	"github.com/jordantipton/golang-restful-webservice/models"
	"github.com/jordantipton/golang-restful-webservice/models/errors"
	"github.com/jordantipton/golang-restful-webservice/services"
)

/*
	Test objects
*/

type mockUserPersister struct {
	mockGetUser    func(userID int) (*models.User, error)
	mockCreateUser func(user *models.User) (*models.User, error)
}

func (m *mockUserPersister) GetUser(userID int) (*models.User, error) {
	if m.mockGetUser != nil {
		return m.mockGetUser(userID)
	}
	return nil, nil
}

func (m *mockUserPersister) CreateUser(user *models.User) (*models.User, error) {
	if m.mockCreateUser != nil {
		return m.mockCreateUser(user)
	}
	return nil, nil
}

/*
	Test functions
*/
func TestGetUserByID(t *testing.T) {
	// Setup
	repositoryUser := &models.User{ID: 1, Name: "Name"}
	mockUserPersister := mockUserPersister{
		mockGetUser: func(userID int) (*models.User, error) {
			if userID == repositoryUser.ID {
				return repositoryUser, nil
			}
			return nil, nil
		},
	}

	usersService := services.UsersService{UsersPersister: &mockUserPersister}

	// Execute
	user, err := usersService.GetUser(1)

	// Assert
	if err != nil {
		t.Errorf("GetUser returned error: %s", err.Error())
	}
	if user == nil {
		t.Errorf("GetUser returned nil")
	}
	if user.ID != repositoryUser.ID {
		t.Errorf("ID, expected: %d, got: %d", repositoryUser.ID, user.ID)
	}
	if user.Name != repositoryUser.Name {
		t.Errorf("ID, expected: %s, got: %s", repositoryUser.Name, user.Name)
	}
}

func TestGetUserByIDNotFound(t *testing.T) {
	// Setup
	mockUserPersister := mockUserPersister{
		mockGetUser: func(userID int) (*models.User, error) {
			return nil, errors.NotFound{Message: "sql: no rows in result set"}
		},
	}

	usersService := services.UsersService{UsersPersister: &mockUserPersister}

	// Execute
	user, err := usersService.GetUser(1)

	// Assert
	if user != nil {
		t.Errorf("Expected user to be nil")
	}
	if err == nil {
		t.Errorf("Expected error to not not be nil")
	}
	if _, ok := err.(errors.NotFound); !ok {
		t.Errorf("Error, expected: NotFound, got: %s", err.Error())
	}
}

func TestGetUserByIDError(t *testing.T) {
	// Setup
	mockUserPersister := mockUserPersister{
		mockGetUser: func(userID int) (*models.User, error) {
			return nil, fmt.Errorf("some error")
		},
	}

	usersService := services.UsersService{UsersPersister: &mockUserPersister}

	// Execute
	_, err := usersService.GetUser(1)

	// Assert
	if err == nil {
		t.Errorf("Expected error to be returned but is nil")
	}
}

func TestCreateUser(t *testing.T) {
	// Setup
	repositoryUser := &models.User{ID: 1, Name: "Name"}
	mockUserPersister := mockUserPersister{
		mockCreateUser: func(*models.User) (*models.User, error) {
			return repositoryUser, nil
		},
	}

	usersService := services.UsersService{UsersPersister: &mockUserPersister}

	// Execute
	user, err := usersService.CreateUser(&models.User{Name: repositoryUser.Name})

	// Assert
	if err != nil {
		t.Errorf("CreateUser returned error: %s", err.Error())
	}
	if user == nil {
		t.Errorf("CreateUser returned nil")
	}
	if user.ID != repositoryUser.ID {
		t.Errorf("ID, expected: %d, got: %d", repositoryUser.ID, user.ID)
	}
	if user.Name != repositoryUser.Name {
		t.Errorf("ID, expected: %s, got: %s", repositoryUser.Name, user.Name)
	}
}

func TestCreateUserNil(t *testing.T) {
	// Setup
	usersService := services.UsersService{UsersPersister: &mockUserPersister{}}

	// Execute
	_, err := usersService.CreateUser(nil)

	// Assert
	if _, ok := err.(errors.InvalidArgument); !ok {
		t.Errorf("Error, expected: InvalidArgument, got: %s", err.Error())
	}
}

func TestCreateUserNoName(t *testing.T) {
	// Setup
	repositoryUser := &models.User{ID: 1, Name: ""}

	usersService := services.UsersService{UsersPersister: &mockUserPersister{}}

	// Execute
	_, err := usersService.CreateUser(&models.User{Name: repositoryUser.Name})

	// Assert
	if _, ok := err.(errors.InvalidArgument); !ok {
		t.Errorf("Error, expected: InvalidArgument, got: %s", err.Error())
	}
}

func TestCreateUserError(t *testing.T) {
	// Setup
	repositoryUser := &models.User{ID: 1, Name: "Name"}
	mockUserPersister := mockUserPersister{
		mockCreateUser: func(*models.User) (*models.User, error) {
			return nil, fmt.Errorf("some error")
		},
	}

	usersService := services.UsersService{UsersPersister: &mockUserPersister}

	// Execute
	_, err := usersService.CreateUser(&models.User{Name: repositoryUser.Name})

	// Assert
	if err == nil {
		t.Errorf("Expected error to be returned but is nil")
	}
}
