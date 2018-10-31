package services

import (
	"fmt"

	"github.com/jordantipton/golang-restful-webservice/models"
	"github.com/jordantipton/golang-restful-webservice/models/errors"
	"github.com/jordantipton/golang-restful-webservice/services/interfaces"
)

type (
	// UsersServicer interface for user services
	UsersServicer interface {
		GetUser(userID int) (*models.User, error)
		CreateUser(user *models.User) (*models.User, error)
	}

	// UsersService providers user information services
	UsersService struct {
		UsersPersister interfaces.UsersPersister
	}
)

// GetUser by ID
func (usersService *UsersService) GetUser(userID int) (*models.User, error) {
	user, err := usersService.UsersPersister.GetUser(userID)
	if err != nil {
		if _, ok := err.(errors.NotFound); ok {
			return nil, errors.NotFound{Message: fmt.Sprintf("User with ID %d not found", userID)}
		}
		return nil, err
	}
	if user == nil {
		return nil, errors.NotFound{Message: fmt.Sprintf("User with ID %d not found", userID)}
	}
	return user, nil
}

// CreateUser and return created user
func (usersService *UsersService) CreateUser(user *models.User) (*models.User, error) {
	if user == nil {
		return nil, errors.InvalidArgument{Message: "User cannot be nil"}
	}
	if user.Name == "" {
		return nil, errors.InvalidArgument{Message: "User name cannot be empty"}
	}
	resultUser, err := usersService.UsersPersister.CreateUser(user)
	if err != nil {
		return nil, err
	}
	return resultUser, nil
}
