package interfaces

import "github.com/jordantipton/golang-restful-webservice/models"

type (
	// UsersPersister interface for user repositories
	UsersPersister interface {
		GetUser(userID int) (*models.User, error)
		CreateUser(user *models.User) (*models.User, error)
	}
)
