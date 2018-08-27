package services_test

import (
	"testing"

	repositoryErrors "github.com/jordantipton/golang-restful-webservice/repositories/errors"
	repositoryModels "github.com/jordantipton/golang-restful-webservice/repositories/models"
	"github.com/jordantipton/golang-restful-webservice/services"
	"github.com/jordantipton/golang-restful-webservice/services/errors"
)

/*
	Test objects
*/

type mockUserPersister struct {
	mockGetUser func(userID int) (*repositoryModels.User, error)
}

func (m *mockUserPersister) GetUser(userID int) (*repositoryModels.User, error) {
	if m.mockGetUser != nil {
		return m.mockGetUser(userID)
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
