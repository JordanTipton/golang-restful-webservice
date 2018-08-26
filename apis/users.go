package apis

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/jordantipton/golang-restful-webservice/services"
)

type (
	UsersResourcer interface {
		GetUser(res http.ResponseWriter, req *http.Request)
	}

	// usersResource defines handlers for the APIs
	UsersResource struct {
		Service services.UsersService
	}
)

// ServeUsersResource sets up the routing of users endpoints and handlers
func RegisterUsersResource(router *chi.Mux, service services.UsersService) {
	r := &UsersResource{service}
	router.Get("/users/{userID}", r.GetUser)
}

func (r *UsersResource) GetUser(res http.ResponseWriter, req *http.Request) {
	//TODO
	json.NewEncoder(res).Encode("")
}
