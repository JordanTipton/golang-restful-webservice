package apis_test

import (
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi"

	"github.com/jordantipton/golang-restful-webservice/apis"
	"github.com/jordantipton/golang-restful-webservice/apis/models"
	"github.com/jordantipton/golang-restful-webservice/services/errors"
	serviceModels "github.com/jordantipton/golang-restful-webservice/services/models"
)

/*
	Test objects
*/

type mockUsersServicer struct {
	mockGetUser func(userID int) (*serviceModels.User, error)
}

func (m *mockUsersServicer) GetUser(userID int) (*serviceModels.User, error) {
	if m.mockGetUser != nil {
		return m.mockGetUser(userID)
	}
	return nil, nil
}

type mockError string

func (e mockError) Error() string { return string(e) }

/*
	Test functions
*/

func TestGetUserByID(t *testing.T) {
	// Setup
	expectedUser := models.User{
		ID:   1,
		Name: "Name",
	}
	serviceUser := &serviceModels.User{ID: expectedUser.ID, Name: expectedUser.Name}
	mockUsersServicer := mockUsersServicer{
		mockGetUser: func(userID int) (*serviceModels.User, error) {
			if userID == serviceUser.ID {
				return serviceUser, nil
			}
			return nil, nil
		},
	}

	r := chi.NewRouter()
	apis.RegisterUsersResource(r, &mockUsersServicer)

	req := httptest.NewRequest("GET", "http://localhost:8080/users/1", nil)
	w := httptest.NewRecorder()

	// Execute
	r.ServeHTTP(w, req)

	// Assert
	if w.Code != 200 {
		t.Errorf("HTTP status code, expected: %d, got: %d", 200, w.Code)
	}

	actualUser := models.User{}
	json.NewDecoder(w.Body).Decode(&actualUser)

	if expectedUser.ID != actualUser.ID {
		t.Errorf("ID, expected: %d, got: %d", expectedUser.ID, actualUser.ID)
	}
	if expectedUser.Name != actualUser.Name {
		t.Errorf("ID, expected: %s, got: %s", expectedUser.Name, actualUser.Name)
	}
}

func TestGetUserByIDBadRequest(t *testing.T) {
	// Setup
	expectedBody := "UserID must be an integer\n"
	mockUsersServicer := mockUsersServicer{}

	r := chi.NewRouter()
	apis.RegisterUsersResource(r, &mockUsersServicer)

	req := httptest.NewRequest("GET", "http://localhost:8080/users/nan", nil)
	w := httptest.NewRecorder()

	// Execute
	r.ServeHTTP(w, req)

	// Assert
	if w.Code != 400 {
		t.Errorf("HTTP status code, expected: %d, got: %d", 400, w.Code)
	}

	body := w.Body.String()
	if body != expectedBody {
		t.Errorf("Response body, expected: %s, got: %s", expectedBody, body)
	}
}

func TestGetUserByIDNotFound(t *testing.T) {
	// Setup
	userID := 1
	expectedBody := fmt.Sprintf("User with ID %d not found\n", userID)
	mockUsersServicer := mockUsersServicer{
		mockGetUser: func(userID int) (*serviceModels.User, error) {
			return nil, errors.NotFound
		},
	}

	r := chi.NewRouter()
	apis.RegisterUsersResource(r, &mockUsersServicer)

	req := httptest.NewRequest("GET", fmt.Sprintf("http://localhost:8080/users/%d", userID), nil)
	w := httptest.NewRecorder()

	// Execute
	r.ServeHTTP(w, req)

	// Assert
	if w.Code != 404 {
		t.Errorf("HTTP status code, expected: %d, got: %d", 404, w.Code)
	}

	body := w.Body.String()
	if body != expectedBody {
		t.Errorf("Response body, expected: %s, got: %s", expectedBody, body)
	}
}

func TestGetUserByIDServerError(t *testing.T) {
	// Setup
	userID := 1
	expectedBody := "some error\n"
	mockUsersServicer := mockUsersServicer{
		mockGetUser: func(userID int) (*serviceModels.User, error) {
			return nil, mockError("some error")
		},
	}

	r := chi.NewRouter()
	apis.RegisterUsersResource(r, &mockUsersServicer)

	req := httptest.NewRequest("GET", fmt.Sprintf("http://localhost:8080/users/%d", userID), nil)
	w := httptest.NewRecorder()

	// Execute
	r.ServeHTTP(w, req)

	// Assert
	if w.Code != 500 {
		t.Errorf("HTTP status code, expected: %d, got: %d", 500, w.Code)
	}

	body := w.Body.String()
	if body != expectedBody {
		t.Errorf("Response body, expected: %s, got: %s", expectedBody, body)
	}
}
