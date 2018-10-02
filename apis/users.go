package apis

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"bitbucket.org/jordantipton/econvote-core/services/errors"
	"github.com/go-chi/chi"
	"github.com/jordantipton/golang-restful-webservice/apis/converters"
	"github.com/jordantipton/golang-restful-webservice/apis/models"
	"github.com/jordantipton/golang-restful-webservice/services"
)

type (
	// UsersResourcer provides an inverface for user resources
	UsersResourcer interface {
		GetUser(res http.ResponseWriter, req *http.Request)
		CreateUser(res http.ResponseWriter, req *http.Request)
	}

	// UsersResource defines handlers for the APIs
	UsersResource struct {
		Service services.UsersServicer
	}
)

// RegisterUsersResource sets up the routing of users endpoints and handlers
func RegisterUsersResource(router *chi.Mux, service services.UsersServicer) {
	r := &UsersResource{service}
	router.Get("/users/{userID}", r.GetUser)
	router.Post("/users", r.CreateUser)
}

// GetUser by ID
func (r *UsersResource) GetUser(res http.ResponseWriter, req *http.Request) {
	userIDString := chi.URLParam(req, "userID")
	userID, err := strconv.Atoi(userIDString)
	if err != nil {
		http.Error(res, "UserID must be an integer", http.StatusBadRequest)
		return
	}
	serviceUser, err := r.Service.GetUser(userID)
	if err != nil {
		if _, ok := err.(errors.NotFound); ok {
			http.Error(res, fmt.Sprintf("User with ID %d not found", userID), http.StatusNotFound)
		} else {
			http.Error(res, err.Error(), http.StatusInternalServerError)
		}
		return
	}
	user := converters.ToUser(serviceUser)
	json.NewEncoder(res).Encode(user)
}

// CreateUser and return result
func (r *UsersResource) CreateUser(res http.ResponseWriter, req *http.Request) {
	var user models.User
	defer req.Body.Close()
	err := json.NewDecoder(req.Body).Decode(&user)
	if err != nil {
		http.Error(res, err.Error(), 400)
		return
	}
	serviceUser, err := r.Service.CreateUser(converters.FromUser(&user))
	if err != nil {
		if _, ok := err.(errors.InvalidArgument); ok {
			http.Error(res, err.Error(), http.StatusBadRequest)
		} else {
			http.Error(res, err.Error(), http.StatusInternalServerError)
		}
		return
	}
	resultUser := converters.ToUser(serviceUser)
	res.WriteHeader(http.StatusCreated)
	json.NewEncoder(res).Encode(resultUser)
}
