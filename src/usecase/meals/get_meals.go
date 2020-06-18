package meals

import (
	"github.com/gin-gonic/gin"
	"go_api/src/repository"
	response2 "go_api/src/schemes/response"
	"go_api/src/types"
	"go_api/src/utils"
	"net/http"
)

// GetMeals godoc
// @Summary get list of meals withing provided date range
// @Tags catering meals
// @Produce json
// @Param startDate query string false "example: 2006-01-02"
// @Param endDate query string false "example: 2006-01-09"
// @Param limit query int false "limit of returned array"
// @Param id path string false "Catering ID"
// @Success 200 {object} response.GetMealsModel "array of meal readings"
// @Failure 400 {object} types.Error "Error"
// @Router /caterings/{id}/meals [get]
func GetMeals(c *gin.Context) {
	var limitQuery types.PaginationQuery
	var dateQuery types.StartEndDateQuery
	var pathUri types.PathId

	if err := utils.RequestBinderUri(&pathUri, c); err != nil {
		return
	}

	_, err := repository.GetCateringByKey("id", pathUri.ID)

	if err != nil {
		if err.Error() == "record not found" {
			utils.CreateError(http.StatusNotFound, err.Error(), c)
			return
		} else {
			utils.CreateError(http.StatusBadRequest, err.Error(), c)
			return
		}
	}

	if err := utils.RequestBinderQuery(&dateQuery, c); err != nil {
		return
	}

	if err := utils.RequestBinderQuery(&limitQuery, c); err != nil {
		return
	}

	result, total, err := repository.GetMealsDB(limitQuery.Limit, dateQuery, pathUri.ID)

	if err != nil {
		utils.CreateError(http.StatusBadRequest, err.Error(), c)
		return
	}

	c.JSON(http.StatusOK, response2.GetMealsModel{
		Items: result,
		Total: total,
	})
}
