package services_test

import (
	"fmt"
	"testing"

	repositoryErrors "github.com/jordantipton/golang-restful-webservice/repositories/errors"
	repositoryModels "github.com/jordantipton/golang-restful-webservice/repositories/models"
	"github.com/jordantipton/golang-restful-webservice/services"
	"github.com/jordantipton/golang-restful-webservice/services/errors"
	"github.com/jordantipton/golang-restful-webservice/services/models"
)

/*
	Test objects
*/

type mockUserPersister struct {
	mockGetUser    func(userID int) (*repositoryModels.User, error)
	mockCreateUser func(user *repositoryModels.User) (*repositoryModels.User, error)
}

func (m *mockUserPersister) GetUser(userID int) (*repositoryModels.User, error) {
	if m.mockGetUser != nil {
		return m.mockGetUser(userID)
	}
	return nil, nil
}

func (m *mockUserPersister) CreateUser(user *repositoryModels.User) (*repositoryModels.User, error) {
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
	repositoryUser := &repositoryModels.User{ID: 1, Name: "Name"}
	mockUserPersister := mockUserPersister{
		mockGetUser: func(userID int) (*repositoryModels.User, error) {
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
		mockGetUser: func(userID int) (*repositoryModels.User, error) {
			return nil, repositoryErrors.NotFound
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
	if err != errors.NotFound {
		t.Errorf("Error, expected: %s, got: %s", errors.NotFound, err.Error())
	}
}

func TestGetUserByIDError(t *testing.T) {
	// Setup
	mockUserPersister := mockUserPersister{
		mockGetUser: func(userID int) (*repositoryModels.User, error) {
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
	repositoryUser := &repositoryModels.User{ID: 1, Name: "Name"}
	mockUserPersister := mockUserPersister{
		mockCreateUser: func(*repositoryModels.User) (*repositoryModels.User, error) {
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
	if err != errors.BadRequest {
		t.Errorf("Error, expected: %s, got: %s", errors.BadRequest, err.Error())
	}
}

func TestCreateUserNoName(t *testing.T) {
	// Setup
	repositoryUser := &repositoryModels.User{ID: 1, Name: ""}

	usersService := services.UsersService{UsersPersister: &mockUserPersister{}}

	// Execute
	_, err := usersService.CreateUser(&models.User{Name: repositoryUser.Name})

	// Assert
	if err != errors.BadRequest {
		t.Errorf("Error, expected: %s, got: %s", errors.BadRequest, err.Error())
	}
}

func TestCreateUserError(t *testing.T) {
	// Setup
	repositoryUser := &repositoryModels.User{ID: 1, Name: "Name"}
	mockUserPersister := mockUserPersister{
		mockCreateUser: func(*repositoryModels.User) (*repositoryModels.User, error) {
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
