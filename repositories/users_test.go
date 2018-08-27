package repositories_test

import (
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jordantipton/golang-restful-webservice/repositories"
	"github.com/jordantipton/golang-restful-webservice/repositories/models"
)

func TestGetUserByID(t *testing.T) {
	// Setup
	userID, userName := 1, "Bob"
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	rows := sqlmock.NewRows([]string{"id", "name"}).AddRow(userID, userName)

	mock.ExpectQuery(fmt.Sprintf("SELECT id, name FROM user WHERE id=%d", userID)).
		WillReturnRows(rows)

	repository := repositories.UsersRepository{DB: db}

	// Execute
	var user *models.User
	user, err = repository.GetUser(userID)

	// Assert
	if mockErr := mock.ExpectationsWereMet(); mockErr != nil {
		t.Errorf("there were unfulfilled expectations: %s", mockErr)
	}
	if err != nil {
		t.Errorf("GetUser returned error: %s", err.Error())
	}
	if user == nil {
		t.Errorf("GetUser returned nil")
	}
	if user.ID != userID {
		t.Errorf("ID, expected: %d, got: %d", userID, user.ID)
	}
	if user.Name != userName {
		t.Errorf("ID, expected: %s, got: %s", userName, user.Name)
	}
}

func TestGetUserByIDError(t *testing.T) {
	// Setup
	userID := 1
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mock.ExpectQuery(fmt.Sprintf("SELECT id, name FROM user WHERE id=%d", userID)).
		WillReturnError(fmt.Errorf("some error"))

	repository := repositories.UsersRepository{DB: db}

	// Execute
	_, err = repository.GetUser(userID)

	// Assert
	if mockErr := mock.ExpectationsWereMet(); mockErr != nil {
		t.Errorf("there were unfulfilled expectations: %s", mockErr)
	}
	if err == nil {
		t.Errorf("Expected error to be returned but is nil")
	}
}
