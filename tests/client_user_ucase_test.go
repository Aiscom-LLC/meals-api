package tests

import (
	"net/http"
	"testing"

	"github.com/Aiscom-LLC/meals-api/api"
	"github.com/Aiscom-LLC/meals-api/api/middleware"
	"github.com/Aiscom-LLC/meals-api/repository"
	"github.com/appleboy/gofight/v2"
	"github.com/buger/jsonparser"
	"github.com/stretchr/testify/assert"
)

func TestAddClientUser(t *testing.T) {
	r := gofight.New()

	var clientRepo = repository.NewClientRepo()
	var userRepo = repository.NewUserRepo()
	userResult, _ := userRepo.GetByKey("email", "meals@aisnovations.com")
	jwt, _, _ := middleware.Passport().TokenGenerator(&middleware.UserID{ID: userResult.ID.String()})
	clientResult, _ := clientRepo.GetByKey("name", "Dymi")
	clientID := clientResult.ID.String()
	email := "testssss@mail.ru"
	var newUserID string

	// Trying to create new user
	// Should be success
	r.POST("/clients/"+clientID+"/users").
		SetCookie(gofight.H{
			"jwt": jwt,
		}).
		SetJSON(gofight.D{
			"email":     email,
			"firstName": "newFirstName",
			"floor":     5,
			"lastName":  "NewLastName",
			"role":      "User",
		}).
		Run(api.SetupRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			data := r.Body.Bytes()
			newUserID, _ = jsonparser.GetString(data, "id")
			assert.Equal(t, http.StatusCreated, r.Code)
		})

	// Trying to create second new user
	// Should be success
	r.POST("/clients/"+clientID+"/users").
		SetCookie(gofight.H{
			"jwt": jwt,
		}).
		SetJSON(gofight.D{
			"email":     "newUserEmail@mail.ru",
			"firstName": "newFirstName",
			"floor":     5,
			"lastName":  "NewLastName",
			"role":      "User",
		}).
		Run(api.SetupRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, http.StatusCreated, r.Code)
		})

	// Trying to create user with already existing email
	// Should return an error
	r.POST("/clients/"+clientID+"/users").
		SetCookie(gofight.H{
			"jwt": jwt,
		}).
		SetJSON(gofight.D{
			"email":     "newUserEmail@mail.ru",
			"firstName": "newFirstName",
			"floor":     5,
			"lastName":  "NewLast_Name",
			"role":      "User",
		}).
		Run(api.SetupRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			data := r.Body.Bytes()
			errorValue, _ := jsonparser.GetString(data, "error")
			assert.Equal(t, http.StatusBadRequest, r.Code)
			assert.Equal(t, "user with that email already exist", errorValue)
		})

	// Trying to create user with invalid email
	// Should return an error
	r.POST("/clients/"+clientID+"/users").
		SetCookie(gofight.H{
			"jwt": jwt,
		}).
		SetJSON(gofight.D{
			"email":     "SecondNewUserEmail.mail.ru",
			"firstName": "newFirstName",
			"floor":     5,
			"lastName":  "NewLast_Name",
			"role":      "User",
		}).
		Run(api.SetupRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			data := r.Body.Bytes()
			errorValue, _ := jsonparser.GetString(data, "error")
			assert.Equal(t, http.StatusBadRequest, r.Code)
			assert.Equal(t, "email is not valid", errorValue)
		})

	// Trying to delete user
	// Should be success
	r.DELETE("/clients/"+clientID+"/users/"+newUserID).
		SetCookie(gofight.H{
			"jwt": jwt,
		}).
		Run(api.SetupRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, http.StatusNoContent, r.Code)
		})

	// Trying to create new user with email which already exist but have status "deleted"
	// Should be success
	r.POST("/clients/"+clientID+"/users").
		SetCookie(gofight.H{
			"jwt": jwt,
		}).
		SetJSON(gofight.D{
			"email":     email,
			"firstName": "newFirstName",
			"floor":     5,
			"lastName":  "NewLast_Name",
			"role":      "User",
		}).
		Run(api.SetupRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, http.StatusCreated, r.Code)
		})
}

func TestGetClientUsers(t *testing.T) {
	r := gofight.New()

	var clientRepo = repository.NewClientRepo()
	var userRepo = repository.NewUserRepo()
	result, _ := userRepo.GetByKey("email", "meals@aisnovations.com")
	jwt, _, _ := middleware.Passport().TokenGenerator(&middleware.UserID{ID: result.ID.String()})
	clientResult, _ := clientRepo.GetByKey("name", "Dymi")
	clientID := clientResult.ID.String()

	// Trying to get list of users
	// Should be success
	r.GET("/clients/"+clientID+"/users?limit=5").
		SetCookie(gofight.H{
			"jwt": jwt,
		}).
		Run(api.SetupRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, http.StatusOK, r.Code)
		})

	// Trying to get user with non-valid catering ID
	// Should return an error
	r.GET("/clients/qwerty/users?limit=5").
		SetCookie(gofight.H{
			"jwt": jwt,
		}).
		Run(api.SetupRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, http.StatusBadRequest, r.Code)
		})
}

func TestDeleteClientUsers(t *testing.T) {
	r := gofight.New()

	var clientRepo = repository.NewClientRepo()
	var userRepo = repository.NewUserRepo()
	result, _ := userRepo.GetByKey("email", "meals@aisnovations.com")
	admin1, _ := userRepo.GetByKey("email", "marianafox@comcubine.com")
	admin2, _ := userRepo.GetByKey("email", "maggietodd@comcubine.com")
	clientAdmin, _ := userRepo.GetByKey("email", "melodybond@comcubine.com")
	adminJWT, _, _ := middleware.Passport().TokenGenerator(&middleware.UserID{ID: clientAdmin.ID.String()})
	jwt, _, _ := middleware.Passport().TokenGenerator(&middleware.UserID{ID: result.ID.String()})
	clientResult, _ := clientRepo.GetByKey("name", "Dymi")
	clientID := clientResult.ID.String()
	var userID string

	// Trying to create new client user
	// Should be success
	r.POST("/clients/"+clientID+"/users").
		SetCookie(gofight.H{
			"jwt": jwt,
		}).
		SetJSON(gofight.D{
			"email":     "newUserEmails@mail.ru",
			"firstName": "newFirstName",
			"floor":     5,
			"lastName":  "NewLastName",
			"role":      "User",
		}).
		Run(api.SetupRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			data := r.Body.Bytes()
			userID, _ = jsonparser.GetString(data, "id")
			assert.Equal(t, http.StatusCreated, r.Code)
		})

	// Trying to delete client user
	// Should be success
	r.DELETE("/clients/"+clientID+"/users/"+userID).
		SetCookie(gofight.H{
			"jwt": jwt,
		}).
		Run(api.SetupRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, http.StatusNoContent, r.Code)
		})

	// Delete client admin user
	// Must be success
	r.DELETE("/clients/"+clientID+"/users/"+admin1.ID.String()).
		SetCookie(gofight.H{
			"jwt": jwt,
		}).
		Run(api.SetupRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, http.StatusNoContent, r.Code)
		})

	r.DELETE("/clients/"+clientID+"/users/"+admin2.ID.String()).
		SetCookie(gofight.H{
			"jwt": adminJWT,
		}).
		Run(api.SetupRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, http.StatusNoContent, r.Code)
		})

	// Trying to delete SUPER USER
	// Should return an error
	r.DELETE("/clients/"+clientID+"/users/"+result.ID.String()).
		SetCookie(gofight.H{
			"jwt": jwt,
		}).
		Run(api.SetupRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			data := r.Body.Bytes()
			errorValue, _ := jsonparser.GetString(data, "error")
			assert.Equal(t, http.StatusBadRequest, r.Code)
			assert.Equal(t, "can't delete yourself", errorValue)
		})
}

func TestUpdateClientUser(t *testing.T) {
	r := gofight.New()

	var clientRepo = repository.NewClientRepo()
	var userRepo = repository.NewUserRepo()
	result, _ := userRepo.GetByKey("email", "meals@aisnovations.com")
	jwt, _, _ := middleware.Passport().TokenGenerator(&middleware.UserID{ID: result.ID.String()})
	clientResult, _ := clientRepo.GetByKey("name", "Dymi")
	clientID := clientResult.ID.String()
	var userID string

	// Trying to create new client user
	// Should be success
	r.POST("/clients/"+clientID+"/users").
		SetCookie(gofight.H{
			"jwt": jwt,
		}).
		SetJSON(gofight.D{
			"email":     "newUserrEmafdsils@mail.ru",
			"firstName": "newFirstName",
			"floor":     5,
			"lastName":  "NewLastName",
			"role":      "User",
		}).
		Run(api.SetupRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			data := r.Body.Bytes()
			userID, _ = jsonparser.GetString(data, "id")
			assert.Equal(t, http.StatusCreated, r.Code)
		})

	// Trying to change name for user
	// Should be success
	r.PUT("/clients/"+clientID+"/users/"+userID).
		SetCookie(gofight.H{
			"jwt": jwt,
		}).
		SetJSON(gofight.D{
			"lastName": "testName",
		}).
		Run(api.SetupRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, http.StatusOK, r.Code)
		})

	// Trying to change email to invalid
	// Should return an error
	r.PUT("/clients/"+clientID+"/users/"+userID).
		SetCookie(gofight.H{
			"jwt": jwt,
		}).
		SetJSON(gofight.D{
			"email": "newwwCoolNameddsaas",
		}).
		Run(api.SetupRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			data := r.Body.Bytes()
			errorValue, _ := jsonparser.GetString(data, "error")
			assert.Equal(t, http.StatusBadRequest, r.Code)
			assert.Equal(t, "email is not valid", errorValue)
		})
}
