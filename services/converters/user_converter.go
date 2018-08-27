package converters

import (
	daoModels "github.com/jordantipton/golang-restful-webservice/repositories/models"
	"github.com/jordantipton/golang-restful-webservice/services/models"
)

// ToUser converts dao User to service User
func ToUser(daoUser *daoModels.User) *models.User {
	return &models.User{ID: daoUser.ID, Name: daoUser.Name}
}
