package repositories

import (
	"database/sql"

	"github.com/jordantipton/golang-restful-webservice/repositories/models"
)

type (
	UsersPersister interface {
		Get(userID int) (*models.User, error)
	}

	UsersRepository struct {
		DB *sql.DB
	}
)

// Get user by ID
func (usersRepository *UsersRepository) Get(userID int) (*models.User, error) {
	//TODO
	//usersRepository.db.
	return &models.User{}, nil
}
