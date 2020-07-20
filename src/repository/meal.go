package repository

import (
	"errors"
	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
	"go_api/src/config"
	"go_api/src/domain"
	"net/http"
	"time"
)

type mealRepo struct{}

func NewMealRepo() *mealRepo {
	return &mealRepo{}
}

// Looks for meal in db, if it doesn't exist returns nil
func (m mealRepo) Find(meal domain.Meal) error {
	if exist := config.DB.
		Where("catering_id = ? AND date = ?", meal.CateringID, meal.Date).
		Find(&meal).RecordNotFound(); !exist {
		return errors.New("this meal already exist")
	}
	return nil
}

// Create meal entity
// returns new meal item and error
func (m mealRepo) Add(meal domain.Meal) (interface{}, error) {
	mealItem := config.DB.Create(&meal)
	if mealItem.Error != nil {
		return nil, mealItem.Error
	}

	return mealItem.Value, nil
}

// Returns list of meals withing provided date range
// Returns list of meals, total items if and error
func (m mealRepo) Get(mealDate time.Time, id string) ([]domain.GetMealDish, uuid.UUID, error, int) {
	var meal domain.Meal
	var result []domain.GetMealDish

	if err := config.DB.
		Where("catering_id = ? AND date = ?", id, mealDate).
		First(&meal).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return []domain.GetMealDish{}, uuid.Nil, errors.New(err.Error()), http.StatusNotFound
		}
		return []domain.GetMealDish{}, uuid.Nil, errors.New(err.Error()), http.StatusBadRequest
	}

	err := config.DB.
		Model(&domain.Category{}).
		Select("categories.name as category_name, categories.id as category_id, d.*").
		Joins("left join dishes d on d.category_id = categories.id").
		Joins("left join meal_dishes md on md.dish_id = d.id").
		Joins("left join meals m on m.id = md.meal_id").
		Where("m.id = ?", meal.ID).
		Scan(&result).
		Error

	return result, meal.ID, err, http.StatusBadRequest
}

func (m mealRepo) GetByKey(key, value string) (domain.Meal, error, int) {
	var meal domain.Meal
	if err := config.DB.
		Where(key+" = ?", value).
		First(&meal).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return domain.Meal{}, errors.New("meal with this date not found"), http.StatusNotFound
		}
		return domain.Meal{}, err, http.StatusBadRequest
	}
	return meal, nil, 0
}