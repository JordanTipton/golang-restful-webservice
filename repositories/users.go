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
		CreateUser(user *models.User) (*models.User, error)
	}

	// UsersRepository represents a repository for user information
	UsersRepository struct {
		DB *sql.DB
	}
)

// GetUser by ID
func (repository *UsersRepository) GetUser(userID int) (*models.User, error) {
	user := models.User{}
	stmt, err := repository.DB.Prepare("SELECT id, name FROM user WHERE id=?")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	err = stmt.QueryRow(userID).Scan(&user.ID, &user.Name)
	if err != nil {
		if err.Error() == errors.SQLNotFound {
			return nil, errors.NotFound{Message: fmt.Sprintf("User with ID %d not found", userID)}
		}
		return nil, err
	}
	return &user, nil
}

// CreateUser in repository and return repository
func (repository *UsersRepository) CreateUser(user *models.User) (*models.User, error) {
	resultUser := models.User{}
	stmtInsert, err := repository.DB.Prepare("INSERT INTO user (name) values(?)")
	if err != nil {
		return nil, err
	}
	defer stmtInsert.Close()
	result, err := stmtInsert.Exec(user.Name)
	if err != nil {
		return nil, err
	}
	stmtSelect, err := repository.DB.Prepare("SELECT id, name FROM user WHERE id=?")
	if err != nil {
		return nil, err
	}
	defer stmtSelect.Close()
	err = stmtSelect.QueryRow(result.LastInsertId()).Scan(&resultUser.ID, &resultUser.Name)
	if err != nil {
		return nil, err
	}
	return &resultUser, nil
}
