package apis_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi"

	"github.com/jordantipton/golang-restful-webservice/apis"
	"github.com/jordantipton/golang-restful-webservice/apis/dtos"
	models "github.com/jordantipton/golang-restful-webservice/models"
	"github.com/jordantipton/golang-restful-webservice/models/errors"
)

/*
	Test objects
*/

type mockUsersServicer struct {
	mockGetUser    func(userID int) (*models.User, error)
	mockCreateUser func(user *models.User) (*models.User, error)
}

func (m *mockUsersServicer) GetUser(userID int) (*models.User, error) {
	if m.mockGetUser != nil {
		return m.mockGetUser(userID)
	}
	return nil, nil
}

func (m *mockUsersServicer) CreateUser(user *models.User) (*models.User, error) {
	if m.mockCreateUser != nil {
		return m.mockCreateUser(user)
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
	expectedUser := dtos.User{
		ID:   1,
		Name: "Name",
	}
	serviceUser := &models.User{ID: expectedUser.ID, Name: expectedUser.Name}
	mockUsersServicer := mockUsersServicer{
		mockGetUser: func(userID int) (*models.User, error) {
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

	actualUser := dtos.User{}
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
		mockGetUser: func(userID int) (*models.User, error) {
			return nil, errors.NotFound{Message: expectedBody}
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
		mockGetUser: func(userID int) (*models.User, error) {
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

func TestCreateUser(t *testing.T) {
	// Setup
	expectedUser := dtos.User{
		ID:   1,
		Name: "Name",
	}
	serviceUser := &models.User{ID: expectedUser.ID, Name: expectedUser.Name}
	mockUsersServicer := mockUsersServicer{
		mockCreateUser: func(user *models.User) (*models.User, error) {
			return serviceUser, nil
		},
	}

	r := chi.NewRouter()
	apis.RegisterUsersResource(r, &mockUsersServicer)

	bodyBytes, _ := json.Marshal(dtos.User{Name: expectedUser.Name})
	req := httptest.NewRequest("POST", "http://localhost:8080/users", bytes.NewReader(bodyBytes))
	w := httptest.NewRecorder()

	// Execute
	r.ServeHTTP(w, req)

	// Assert
	if w.Code != 201 {
		t.Errorf("HTTP status code, expected: %d, got: %d", 201, w.Code)
	}

	actualUser := dtos.User{}
	json.NewDecoder(w.Body).Decode(&actualUser)

	if expectedUser.ID != actualUser.ID {
		t.Errorf("ID, expected: %d, got: %d", expectedUser.ID, actualUser.ID)
	}
	if expectedUser.Name != actualUser.Name {
		t.Errorf("ID, expected: %s, got: %s", expectedUser.Name, actualUser.Name)
	}
}

func TestCreateUserError(t *testing.T) {
	// Setup
	expectedUser := dtos.User{
		ID:   1,
		Name: "Name",
	}
	expectedBody := "some error"
	mockUsersServicer := mockUsersServicer{
		mockCreateUser: func(user *models.User) (*models.User, error) {
			return nil, fmt.Errorf(expectedBody)
		},
	}

	r := chi.NewRouter()
	apis.RegisterUsersResource(r, &mockUsersServicer)

	bodyBytes, _ := json.Marshal(dtos.User{Name: expectedUser.Name})
	req := httptest.NewRequest("POST", "http://localhost:8080/users", bytes.NewReader(bodyBytes))
	w := httptest.NewRecorder()

	// Execute
	r.ServeHTTP(w, req)

	// Assert
	if w.Code != 500 {
		t.Errorf("HTTP status code, expected: %d, got: %d", 500, w.Code)
	}
	body := w.Body.String()
	if body != expectedBody+"\n" {
		t.Errorf("Response body, expected: %s, got: %s", expectedBody, body)
	}
}

func TestCreateUserNilBodyBadRequest(t *testing.T) {
	// Setup
	mockUsersServicer := mockUsersServicer{
		mockCreateUser: func(user *models.User) (*models.User, error) {
			return nil, errors.InvalidArgument{Message: "Missing body"}
		},
	}

	r := chi.NewRouter()
	apis.RegisterUsersResource(r, &mockUsersServicer)

	req := httptest.NewRequest("POST", "http://localhost:8080/users", nil)
	w := httptest.NewRecorder()

	// Execute
	r.ServeHTTP(w, req)

	// Assert
	if w.Code != 400 {
		t.Errorf("HTTP status code, expected: %d, got: %d", 400, w.Code)
	}
}

func TestCreateUserNoNameBadRequest(t *testing.T) {
	// Setup
	expectedBody := "User must have a name"
	mockUsersServicer := mockUsersServicer{
		mockCreateUser: func(user *models.User) (*models.User, error) {
			return nil, errors.InvalidArgument{Message: expectedBody}
		},
	}

	r := chi.NewRouter()
	apis.RegisterUsersResource(r, &mockUsersServicer)

	bodyBytes, _ := json.Marshal(dtos.User{})
	req := httptest.NewRequest("POST", "http://localhost:8080/users", bytes.NewReader(bodyBytes))
	w := httptest.NewRecorder()

	// Execute
	r.ServeHTTP(w, req)

	// Assert
	if w.Code != 400 {
		t.Errorf("HTTP status code, expected: %d, got: %d", 400, w.Code)
	}
	body := w.Body.String()
	if body != expectedBody+"\n" {
		t.Errorf("Response body, expected: %s, got: %s", expectedBody, body)
	}
}
