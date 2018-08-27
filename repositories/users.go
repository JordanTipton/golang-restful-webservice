package repositories

import (
	"database/sql"
	"fmt"

	"github.com/jordantipton/golang-restful-webservice/repositories/errors"
	"github.com/jordantipton/golang-restful-webservice/repositories/models"
)

type (
	// UsersPersister interface for user repositories
	UsersPersister interface {
		GetUser(userID int) (*models.User, error)
	}

	// UsersRepository represents a repository for user information
	UsersRepository struct {
		DB *sql.DB
	}
)

// GetUser by ID
func (repository *UsersRepository) GetUser(userID int) (*models.User, error) {
	user := models.User{}
	statement := fmt.Sprintf("SELECT id, name FROM user WHERE id=%d", userID)
	err := repository.DB.QueryRow(statement).Scan(&user.ID, &user.Name)
	if err != nil {
		if err.Error() == errors.NotFound.Error() {
			return nil, errors.NotFound
		}
		return nil, err
	}
	return &user, nil
}
