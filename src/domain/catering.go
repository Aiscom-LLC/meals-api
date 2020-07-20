package domain

import (
	"github.com/gin-gonic/gin"
	"go_api/src/types"
)

// Catering model
type Catering struct {
	Base
	Name string `gorm:"type:varchar(30);not null" json:"name,omitempty" binding:"required"`
} //@name CateringsResponse

type CateringUsecase interface {
	Get(c *gin.Context)
	Add(c *gin.Context)
	Update(c *gin.Context)
	Delete(c *gin.Context)
}

type CateringRepository interface {
	Get(query types.PaginationQuery) ([]Catering, int, error)
	Add(catering Catering) error
	Update(id string, catering Catering) (error, int)
	Delete(id string) error
	GetByKey(key, value string) (Catering, error)
}