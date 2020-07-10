package request

import "time"

type AddMeal struct {
	Date time.Time `json:"date" binding:"required" example:"2020-06-20T00:00:00Z"`
} // @name AddMealRequest
