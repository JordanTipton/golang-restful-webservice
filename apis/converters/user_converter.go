package converters

import (
	"github.com/jordantipton/golang-restful-webservice/apis/models"
	serviceModels "github.com/jordantipton/golang-restful-webservice/services/models"
)

// ToUser converts service User to api User
func ToUser(serviceUser *serviceModels.User) *models.User {
	return &models.User{ID: serviceUser.ID, Name: serviceUser.Name}
}

// FromUser converts api User to service User
func FromUser(apiUser *models.User) *serviceModels.User {
	return &serviceModels.User{ID: apiUser.ID, Name: apiUser.Name}
}
