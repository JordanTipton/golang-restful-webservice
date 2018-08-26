package services

import (
	"github.com/jordantipton/golang-restful-webservice/repositories"
	"github.com/jordantipton/golang-restful-webservice/services/models"
)

type (
	UsersServicer interface {
		Get(userID int) (*models.User, error)
	}

	UsersService struct {
		UsersPersister repositories.UsersPersister
	}
)

func (usersService *UsersService) Get(userID int) (*models.User, error) {
	//TODO
	return &models.User{}, nil
}
