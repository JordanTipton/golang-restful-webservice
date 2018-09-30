package services

import (
	"github.com/jordantipton/golang-restful-webservice/repositories"
	repositoryErrors "github.com/jordantipton/golang-restful-webservice/repositories/errors"
	"github.com/jordantipton/golang-restful-webservice/services/converters"
	"github.com/jordantipton/golang-restful-webservice/services/errors"
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
		if err == repositoryErrors.NotFound {
			return nil, errors.NotFound
		}
		return nil, err
	}
	if daoUser == nil {
		return nil, errors.NotFound
	}
	user := converters.ToUser(daoUser)
	return user, nil
}

// CreateUser and return created user
func (usersService *UsersService) CreateUser(user *models.User) (*models.User, error) {
	if user == nil || user.Name == "" {
		return nil, errors.BadRequest
	}
	daoUser, err := usersService.UsersPersister.CreateUser(converters.FromUser(user))
	if err != nil {
		return nil, err
	}
	resultUser := converters.ToUser(daoUser)
	return resultUser, nil
}
