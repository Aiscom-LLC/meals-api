package request

import "time"

type AddCategory struct {
	Name      string     `json:"name" example:"закуски" binding:"required"`
	DeletedAt *time.Time `json:"deletedAt" example:"2020-06-29T00:00:00Z"`
} //@name AddCategoryRequest