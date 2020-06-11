package catering

import (
	"github.com/gin-gonic/gin"
	"go_api/src/repository/catering"
	"go_api/src/types"
	"net/http"
)

// GetCatering godoc
// @Summary Returns list of caterings
// @Tags catering
// @Produce json
// @Param id path string true "Catering ID"
// @Success 200 {object} models.Catering "List of caterings"
// @Failure 400 {object} types.Error "Error"
// @Failure 404 {object} types.Error "Error"
// @Router /caterings/{id} [get]
func GetCatering(c *gin.Context) {
	var path types.PathId
	if err := c.ShouldBindUri(&path); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"code":  http.StatusBadRequest,
			"error": err.Error(),
		})
		return
	}
	result, err := catering.GetCateringByKey("id", path.ID)
	if err != nil {
		if err.Error() == "record not found" {
			c.JSON(http.StatusNotFound, gin.H{
				"code":  http.StatusNotFound,
				"error": err.Error(),
			})
			return
		} else {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":  http.StatusBadRequest,
				"error": err.Error(),
			})
			return
		}
	}
	c.JSON(http.StatusOK, result)
}