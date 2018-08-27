package apis

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/jordantipton/golang-restful-webservice/apis/converters"
	"github.com/jordantipton/golang-restful-webservice/services"
	"github.com/jordantipton/golang-restful-webservice/services/errors"
)

type (
	// UsersResourcer provides an inverface for user resources
	UsersResourcer interface {
		GetUser(res http.ResponseWriter, req *http.Request)
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
		if err == errors.NotFound {
			http.Error(res, fmt.Sprintf("User with ID %d not found", userID), http.StatusNotFound)
		} else {
			http.Error(res, err.Error(), http.StatusInternalServerError)
		}
		return
	}
	user := converters.ToUser(serviceUser)
	json.NewEncoder(res).Encode(user)
}
