package domain

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
)

// MealBase struct without deletedAt
type MealBase struct {
	ID        uuid.UUID  `gorm:"url:uuid;" json:"id"`
	DeletedAt *time.Time `json:"-"`
	UpdatedAt time.Time  `json:"-"`
}

// BeforeCreate func which generates uuid v4 for each inserted row
func (base *MealBase) BeforeCreate(scope *gorm.Scope) error {
	uuidv4 := uuid.NewV4()
	return scope.SetColumn("ID", uuidv4)
}

// Meal struct for DB
type Meal struct {
	MealBase
	CreatedAt  time.Time `json:"createdAt"`
	Date       time.Time `json:"date,omitempty" binding:"required"`
	CateringID uuid.UUID `json:"-"`
	ClientID   uuid.UUID `json:"-"`
	MealID     uuid.UUID `json:"mealId"`
	Version    string    `json:"version"`
	Person     string    `json:"person"`
} // @name MealsResponse

// MealUsecase is meal interface for usecase
type MealUsecase interface {
	Add(c *gin.Context)
	Get(c *gin.Context)
	Update(c *gin.Context)
}

// MealRepository is meal interface for repository
type MealRepository interface {
	Find(meal *Meal) error
	Add(meal *Meal) error
	/* TODO fix cycle imports
	Get(mealDate time.Time, id, clientID string) ([]response.GetMeal, int, error)
	*/
	GetByKey(key, value string) (Meal, int, error)
}
