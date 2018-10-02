package services

import (
	"fmt"

	"bitbucket.org/jordantipton/econvote-core/services/errors"
	"github.com/jordantipton/golang-restful-webservice/repositories"
	repositoryErrors "github.com/jordantipton/golang-restful-webservice/repositories/errors"
	"github.com/jordantipton/golang-restful-webservice/services/converters"
	"github.com/jordantipton/golang-restful-webservice/services/models"
)

type (
	// UsersServicer interface for user services
	UsersServicer interface {
		GetUser(userID int) (*models.User, error)
		CreateUser(user *models.User) (*models.User, error)
	}

	// UsersService providers user information services
	UsersService struct {
		UsersPersister repositories.UsersPersister
	}
)

// GetUser by ID
func (usersService *UsersService) GetUser(userID int) (*models.User, error) {
	daoUser, err := usersService.UsersPersister.GetUser(userID)
	if err != nil {
		if _, ok := err.(repositoryErrors.NotFound); ok {
			return nil, errors.NotFound{Message: fmt.Sprintf("User with ID %d not found", userID)}
		}
		return nil, err
	}
	if daoUser == nil {
		return nil, errors.NotFound{Message: fmt.Sprintf("User with ID %d not found", userID)}
	}
	user := converters.ToUser(daoUser)
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
	daoUser, err := usersService.UsersPersister.CreateUser(converters.FromUser(user))
	if err != nil {
		return nil, err
	}
	resultUser := converters.ToUser(daoUser)
	return resultUser, nil
}
