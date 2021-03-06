package tests

import (
	"net/http"
	"testing"

	"github.com/Aiscom-LLC/meals-api/api"
	"github.com/Aiscom-LLC/meals-api/api/middleware"
	"github.com/Aiscom-LLC/meals-api/repository"

	"github.com/appleboy/gofight/v2"
	"github.com/buger/jsonparser"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
)

func TestAddDish(t *testing.T) {
	r := gofight.New()

	categoryRepo := repository.NewCategoryRepo()
	userRepo := repository.NewUserRepo()
	cateringRepo := repository.NewCateringRepo()
	userResult, _ := userRepo.GetByKey("email", "meals@aisnovations.com")
	jwt, _, _ := middleware.Passport().TokenGenerator(&middleware.UserID{ID: userResult.ID.String()})

	cateringResult, _ := cateringRepo.GetByKey("name", "Twiist")
	cateringID := cateringResult.ID.String()

	categoryResult, _ := categoryRepo.GetByKey("name", "гарнир", cateringID)

	// Trying to add dish to non-existing catering
	// Should throw an error
	fakeID := uuid.NewV4()
	r.POST("/caterings/"+fakeID.String()+"/dishes").
		SetCookie(gofight.H{
			"jwt": jwt,
		}).
		SetJSON(gofight.D{
			"desc":   "Очень вкусный",
			"name":   "тест",
			"price":  120,
			"weight": 250,
		}).
		Run(api.SetupRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, http.StatusBadRequest, r.Code)
		})

	// Trying to add dish to non-existing dish category
	// Should throw an error
	r.POST("/caterings/"+cateringID+"/dishes").
		SetCookie(gofight.H{
			"jwt": jwt,
		}).
		SetJSON(gofight.D{
			"desc":   "Очень вкусный",
			"name":   "тест",
			"price":  120,
			"weight": 250,
		}).
		Run(api.SetupRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, http.StatusBadRequest, r.Code)
		})

	// Trying to create new dish
	// Should be success
	r.POST("/caterings/"+cateringID+"/dishes").
		SetCookie(gofight.H{
			"jwt": jwt,
		}).
		SetJSON(gofight.D{
			"categoryID": categoryResult.ID,
			"desc":       "Очень вкусный",
			"name":       "тест",
			"price":      120,
			"weight":     250,
		}).
		Run(api.SetupRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, http.StatusOK, r.Code)
		})
}

func TestGetDishes(t *testing.T) {
	r := gofight.New()

	categoryRepo := repository.NewCategoryRepo()
	userRepo := repository.NewUserRepo()
	cateringRepo := repository.NewCateringRepo()
	userResult, _ := userRepo.GetByKey("email", "meals@aisnovations.com")
	jwt, _, _ := middleware.Passport().TokenGenerator(&middleware.UserID{ID: userResult.ID.String()})

	cateringResult, _ := cateringRepo.GetByKey("name", "Twiist")
	cateringID := cateringResult.ID.String()

	categoryResult, _ := categoryRepo.GetByKey("name", "супы", cateringID)
	categoryID := categoryResult.ID.String()

	fakeID := uuid.NewV4()

	// Trying to get dishes with non-existing catering ID
	// Should throw an error
	r.GET("/caterings/"+fakeID.String()+"/dishes?categoryID="+categoryID).
		SetCookie(gofight.H{
			"jwt": jwt,
		}).
		Run(api.SetupRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			data := r.Body.Bytes()
			errorValue, _ := jsonparser.GetString(data, "error")
			assert.Equal(t, http.StatusNotFound, r.Code)
			assert.Equal(t, "catering with that ID doesn't exist", errorValue)
		})

	// Trying to get dishes with non-existing category ID
	// Should throw an error
	r.GET("/caterings/"+cateringID+"/dishes?categoryID="+fakeID.String()).
		SetCookie(gofight.H{
			"jwt": jwt,
		}).
		Run(api.SetupRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			data := r.Body.Bytes()
			errorValue, _ := jsonparser.GetString(data, "error")
			assert.Equal(t, http.StatusNotFound, r.Code)
			assert.Equal(t, "category with that ID doesn't exist", errorValue)
		})

	// Trying to get dishes with all valid values
	// Should be success
	r.GET("/caterings/"+cateringID+"/dishes?categoryID="+categoryResult.ID.String()).
		SetCookie(gofight.H{
			"jwt": jwt,
		}).
		Run(api.SetupRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, http.StatusOK, r.Code)
		})
}

func TestUpdateDish(t *testing.T) {
	r := gofight.New()

	categoryRepo := repository.NewCategoryRepo()
	userRepo := repository.NewUserRepo()
	dishRepo := repository.NewDishRepo()
	cateringRepo := repository.NewCateringRepo()
	userResult, _ := userRepo.GetByKey("email", "meals@aisnovations.com")
	jwt, _, _ := middleware.Passport().TokenGenerator(&middleware.UserID{ID: userResult.ID.String()})

	cateringResult, _ := cateringRepo.GetByKey("name", "Twiist")
	cateringID := cateringResult.ID.String()

	categoryResult, _ := categoryRepo.GetByKey("name", "супы", cateringID)
	categoryID := categoryResult.ID.String()

	dishResult, _, _ := dishRepo.GetByKey("name", "борщ", cateringID, categoryID)
	dishID := dishResult.ID.String()

	fakeID := uuid.NewV4()

	// Trying to update dish for non-existing catering
	// Should throw an error
	r.PUT("/caterings/"+fakeID.String()+"/dishes/"+dishID).
		SetCookie(gofight.H{
			"jwt": jwt,
		}).
		SetJSON(gofight.D{
			"categoryID": categoryID,
			"desc":       "Очень острый",
			"name":       "супер доширак",
			"price":      120,
			"weight":     250,
		}).
		Run(api.SetupRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			data := r.Body.Bytes()
			errorValue, _ := jsonparser.GetString(data, "error")
			assert.Equal(t, http.StatusNotFound, r.Code)
			assert.Equal(t, "catering not found", errorValue)
		})

	// Trying to update dish with non-existing dish category id
	// Should throw an error
	r.PUT("/caterings/"+cateringID+"/dishes/"+dishID).
		SetCookie(gofight.H{
			"jwt": jwt,
		}).
		SetJSON(gofight.D{
			"categoryID": fakeID,
			"desc":       "Очень острый",
			"name":       "супер доширак",
			"price":      120,
			"weight":     250,
		}).
		Run(api.SetupRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			data := r.Body.Bytes()
			errorValue, _ := jsonparser.GetString(data, "error")
			assert.Equal(t, http.StatusNotFound, r.Code)
			assert.Equal(t, "dish category not found", errorValue)
		})

	// Trying to update dish with non-existing dish id
	// Should throw an error
	r.PUT("/caterings/"+cateringID+"/dishes/"+fakeID.String()).
		SetCookie(gofight.H{
			"jwt": jwt,
		}).
		SetJSON(gofight.D{
			"categoryID": categoryID,
			"desc":       "Очень острый",
			"name":       "супер доширак",
			"price":      120,
			"weight":     250,
		}).
		Run(api.SetupRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			data := r.Body.Bytes()
			errorValue, _ := jsonparser.GetString(data, "error")
			assert.Equal(t, http.StatusNotFound, r.Code)
			assert.Equal(t, "dish not found", errorValue)
		})

	// Trying to update dish with all valid values
	// Should be success
	r.PUT("/caterings/"+cateringID+"/dishes/"+dishID).
		SetCookie(gofight.H{
			"jwt": jwt,
		}).
		SetJSON(gofight.D{
			"categoryID": categoryID,
			"desc":       "Очень острый",
			"name":       "супер доширак",
			"price":      120,
			"weight":     250,
		}).
		Run(api.SetupRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, http.StatusNoContent, r.Code)
		})
}

func TestDeleteDish(t *testing.T) {
	r := gofight.New()

	categoryRepo := repository.NewCategoryRepo()
	userRepo := repository.NewUserRepo()
	dishRepo := repository.NewDishRepo()
	cateringRepo := repository.NewCateringRepo()
	userResult, _ := userRepo.GetByKey("email", "meals@aisnovations.com")
	jwt, _, _ := middleware.Passport().TokenGenerator(&middleware.UserID{ID: userResult.ID.String()})

	cateringResult, _ := cateringRepo.GetByKey("name", "Twiist")
	cateringID := cateringResult.ID.String()

	categoryResult, _ := categoryRepo.GetByKey("name", "супы", cateringID)
	categoryID := categoryResult.ID.String()

	dishResult, _, _ := dishRepo.GetByKey("name", "доширак", cateringID, categoryID)

	fakeID := uuid.NewV4()

	// Trying to delete non-existing dish
	// Should throw an error
	r.DELETE("/caterings/"+cateringID+"/dishes/"+fakeID.String()).
		SetCookie(gofight.H{
			"jwt": jwt,
		}).
		Run(api.SetupRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			data := r.Body.Bytes()
			errorValue, _ := jsonparser.GetString(data, "error")
			assert.Equal(t, http.StatusNotFound, r.Code)
			assert.Equal(t, "dish not found", errorValue)
		})

	// Trying to delete existing dish
	// Should be success
	r.DELETE("/caterings/"+cateringID+"/dishes/"+dishResult.ID.String()).
		SetCookie(gofight.H{
			"jwt": jwt,
		}).
		Run(api.SetupRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, http.StatusNoContent, r.Code)
		})

	// Trying to delete soft deleted dish
	// Should throw an error
	r.DELETE("/caterings/"+cateringID+"/dishes/"+dishResult.ID.String()).
		SetCookie(gofight.H{
			"jwt": jwt,
		}).
		Run(api.SetupRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			data := r.Body.Bytes()
			errorValue, _ := jsonparser.GetString(data, "error")
			assert.Equal(t, http.StatusNotFound, r.Code)
			assert.Equal(t, "dish not found", errorValue)
		})
}
