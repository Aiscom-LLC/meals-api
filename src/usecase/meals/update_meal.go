package meals

import (
	"github.com/gin-gonic/gin"
	"go_api/src/models"
	"go_api/src/repository/meal"
	meal2 "go_api/src/schemes/request/meal"
	"go_api/src/types"
	"go_api/src/utils"
	"net/http"
)

// UpdateMeal godoc
// @Summary Returns updated meal
// @Produce json
// @Accept json
// @Tags meals
// @Param id path string true "Meal ID"
// @Param body body meal.UpdateMealRequest false "Meal date"
// @Success 200 {object} models.Meal "Meal"
// @Failure 400 {object} types.Error "Error"
// @Router /meals/{id} [put]
func UpdateMeal(c *gin.Context) {
	var path types.PathId
	var body meal2.UpdateMealRequest

	if err := utils.RequestBinderUri(&path, c); err != nil {
		return
	}

	if err := utils.RequestBinderBody(&body, c); err != nil {
		return
	}

	result, err := meal.UpdateMealDB(path.ID, body)

	if err != nil {
		utils.CreateError(http.StatusBadRequest, err.Error(), c)
		return
	}

	if result.RowsAffected == 0 {
		if result.Error != nil {
			utils.CreateError(http.StatusBadRequest, result.Error.Error(), c)
			return
		}

		utils.CreateError(http.StatusNotFound, "meal not found", c)
		return
	}

	mealDate := result.Value.(*models.Meal).Date

	c.JSON(http.StatusOK, gin.H{
		"id":   path.ID,
		"date": mealDate,
	})

}
