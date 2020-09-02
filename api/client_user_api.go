package api

import (
	"errors"
	"github.com/Aiscom-LLC/meals-api/domain"
	"github.com/Aiscom-LLC/meals-api/mailer"
	"github.com/Aiscom-LLC/meals-api/repository"
	"github.com/Aiscom-LLC/meals-api/schemes/request"
	"github.com/Aiscom-LLC/meals-api/services"
	"github.com/Aiscom-LLC/meals-api/types"
	"github.com/Aiscom-LLC/meals-api/utils"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"net/http"
)

// ClientUser struct
type ClientUser struct{}

// NewClientUser returns pointer to client
// user struct with all methods
func NewClientUser() *ClientUser {
	return &ClientUser{}
}

var userRepo = repository.NewUserRepo()
var clientUserRepo = repository.NewClientUserRepo()
var clientUserService = services.NewClientUser()

// Add creates user for client
// @Summary Returns error or 201 status code if success
// @Produce json
// @Accept json
// @Tags clients users
// @Param id path string false "Client ID"
// @Param body body request.ClientUser false "Client user"
// @Success 201 {object} response.UserResponse false "Client user"
// @Failure 400 {object} types.Error "Error"
// @Router /clients/{id}/users [post]
func (cu ClientUser) Add(c *gin.Context) {
	var path types.PathID
	var body request.ClientUser
	var user domain.User

	if err := utils.RequestBinderURI(&path, c); err != nil {

		return
	}

	if err := utils.RequestBinderBody(&body, c); err != nil {

		return
	}

	if err := copier.Copy(&user, &body); err != nil {

		utils.CreateError(http.StatusBadRequest, err, c)
		return
	}

	if ok := utils.IsEmailValid(user.Email); !ok {

		utils.CreateError(http.StatusBadRequest, errors.New("email is not valid"), c)
		return
	}

	userClientCatering, code, err, userErr, password := clientUserService.Add(path, body, user)

	if err != nil {

		utils.CreateError(code, err, c)
		return
	}
	if userErr != nil {

		utils.CreateError(code, userErr, c)
		return
	}

	// nolint:errcheck
	go mailer.SendEmail(user, password)
	c.JSON(http.StatusCreated, userClientCatering)
}

// Get return list of client users
// @Summary Returns list of client users
// @Tags clients users
// @Produce json
// @Param id path string false "Client ID"
// @Param limit query int false "used for pagination"
// @Param page query int false "used for pagination"
// @Param q query string false "used query search"
// @Param role query string false "used for role sort"
// @Param status query string false "used for status sort"
// @Success 200 {array} response.UserResponse "List of client users"
// @Failure 400 {object} types.Error "Error"
// @Failure 404 {object} types.Error "Error"
// @Router /clients/{id}/users [get]
func (cu ClientUser) Get(c *gin.Context) { //nolint:dupl
	var path types.PathID
	var paginationQuery types.PaginationQuery
	var filterQuery types.UserFilterQuery

	if err := utils.RequestBinderURI(&path, c); err != nil {

		return
	}

	if err := utils.RequestBinderQuery(&paginationQuery, c); err != nil {

		return
	}

	if err := utils.RequestBinderQuery(&filterQuery, c); err != nil {

		return
	}

	ctxUser, _ := c.Get("user")
	ctxUserRole := ctxUser.(domain.User).Role

	users, total, code, err := clientUserRepo.Get(path.ID, ctxUserRole, paginationQuery, filterQuery)

	if err != nil {

		utils.CreateError(code, err, c)
		return
	}

	if paginationQuery.Page == 0 {
		paginationQuery.Page = 1
	}

	c.JSON(http.StatusOK, gin.H{
		"items": users,
		"total": total,
		"page":  paginationQuery.Page,
	})
}

// Delete deletes user of client
// @Summary Returns error or 204 status code if success
// @Produce json
// @Accept json
// @Tags clients users
// @Param id path string false "Client ID"
// @Param userId path string false "User ID"
// @Success 204 "Successfully deleted"
// @Failure 400 {object} types.Error "Error"
// @Failure 404 {object} types.Error "Error"
// @Router /clients/{id}/users/{userId} [delete]
func (cu ClientUser) Delete(c *gin.Context) {
	var path types.PathUser
	var user domain.User

	if err := utils.RequestBinderURI(&path, c); err != nil {

		return
	}

	delUser, _ := c.Get("user")
	userRole := delUser.(domain.User).Role
	userID := delUser.(domain.User).ID.String()
	code, err := clientUserService.Delete(path, user, userRole, userID)

	if err != nil {

		utils.CreateError(code, err, c)
		return
	}

	c.Status(http.StatusNoContent)
}

// Update updates user of client
// @Summary Returns error or 200 status code if success
// @Produce json
// @Accept json
// @Tags clients users
// @Param id path string false "Client ID"
// @Param userId path string false "User ID"
// @Param body body request.ClientUserUpdate false "Client user"
// @Success 200 {object} response.UserResponse false "Client user"
// @Failure 400 {object} types.Error "Error"
// @Failure 404 {object} types.Error "Error"
// @Router /clients/{id}/users/{userId} [put]
func (cu ClientUser) Update(c *gin.Context) {
	var path types.PathUser
	var body request.ClientUserUpdate
	var user domain.User

	if err := utils.RequestBinderURI(&path, c); err != nil {

		return
	}

	if err := utils.RequestBinderBody(&body, c); err != nil {

		return
	}

	code, err := clientUserService.Update(path, body, user)

	if err != nil {

		utils.CreateError(code, err, c)
		return
	}
	updatedUser, _ := userRepo.GetByID(path.UserID)
	c.JSON(http.StatusOK, updatedUser)

}