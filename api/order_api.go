package api

import (
	"github.com/Aiscom-LLC/meals-api/api/middleware"
	"github.com/Aiscom-LLC/meals-api/schemes/request"
	"github.com/Aiscom-LLC/meals-api/schemes/response"
	"github.com/Aiscom-LLC/meals-api/services"
	"github.com/Aiscom-LLC/meals-api/types"
	"github.com/Aiscom-LLC/meals-api/utils"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"net/http"
)

// Order struct
type Order struct{}

// NewOrder return pointer to order struct
// with all methods
func NewOrder() *Order {
	return &Order{}
}

// Add creates order for client user
// @Summary Returns error or 201 status code if success
// @Produce json
// @Accept json
// @Tags users orders
// @Param id path string false "User ID"
// @Param date query string true "Date query in YYYY-MM-DDT00:00:00Z format"
// @Param body body request.OrderRequest false "User order"
// @Success 201 {object} response.UserOrder false "Order for user"
// @Failure 400 {object} types.Error "Error"
// @Router /users/{id}/orders [post]
func (o Order) Add(c *gin.Context) {
	var query types.DateQuery
	var order request.OrderRequest
	var path types.PathID

	claims, err := middleware.Passport().GetClaimsFromJWT(c)

	if err != nil {
		return
	}

	if err := utils.RequestBinderQuery(&query, c); err != nil {
		return
	}

	if err := utils.RequestBinderBody(&order, c); err != nil {
		return
	}

	if err := utils.RequestBinderURI(&path, c); err != nil {
		return
	}

	userOrder, code, err := services.NewOrderService().Add(query.Date, order, jwt.MapClaims(claims))

	if err != nil {
		utils.CreateError(code, err, c)
		return
	}

	c.JSON(http.StatusCreated, response.UserOrder{
		Items:   userOrder.Items,
		Status:  userOrder.Status,
		Total:   userOrder.Total,
		OrderID: userOrder.OrderID,
	})
}

// CancelOrder changes status of order to canceled
// @Summary Returns error or 204 status code if success
// @Produce json
// @Accept json
// @Tags users orders
// @Param id path string false "User ID"
// @Param orderId path string false "Order ID"
// @Success 204 "Successfully canceled"
// @Failure 400 {object} types.Error "Error"
// @Failure 404 {object} types.Error "Error"
// @Router /users/{id}/orders/{orderId} [delete]
func (o Order) CancelOrder(c *gin.Context) {
	var path types.PathOrder

	if err := utils.RequestBinderURI(&path, c); err != nil {
		return
	}

	code, err := services.NewOrderService().CancelOrderService(path)

	if err != nil {
		utils.CreateError(code, err, c)
		return
	}

	c.Status(http.StatusNoContent)
}

// GetUserOrder returns orders of provided user
// @Summary returns orders of provided user
// @Tags users orders
// @Produce json
// @Param id path string true "User ID"
// @Param date query string true "Date query in YYYY-MM-DDT00:00:00Z format"
// @Success 200 {array} response.UserOrder false "Orders for user"
// @Failure 400 {object} types.Error "Error"
// @Failure 404 {object} types.Error "Not Found"
// @Router /users/{id}/orders [get]
func (o Order) GetUserOrder(c *gin.Context) {
	var path types.PathID
	var query types.DateQuery

	if err := utils.RequestBinderURI(&path, c); err != nil {
		return
	}

	if err := utils.RequestBinderQuery(&query, c); err != nil {
		return
	}

	userOrder, code, err := services.NewOrderService().GetUserOrderService(path, query)

	if err != nil {
		utils.CreateError(code, err, c)
		return
	}

	c.JSON(http.StatusOK, userOrder)
}

// GetClientOrders returns orders of provided client
// @Summary returns orders of provided client
// @Tags clients orders
// @Produce json
// @Param id path string true "Client ID"
// @Param date query string true "Date query in YYYY-MM-DDT00:00:00Z format"
// @Success 200 {object} response.SummaryOrdersResponse false "Orders for clients"
// @Failure 400 {object} types.Error "Error"
// @Failure 404 {object} types.Error "Not Found"
// @Router /clients/{id}/orders [get]
func (o Order) GetClientOrders(c *gin.Context) {
	var path types.PathID
	var query types.DateQuery

	if err := utils.RequestBinderURI(&path, c); err != nil {
		return
	}

	if err := utils.RequestBinderQuery(&query, c); err != nil {
		return
	}

	client := types.CompanyTypesEnum.Client

	result, code, err := services.NewOrderService().GetClientOrderService(path, query, client)

	if err != nil {
		utils.CreateError(code, err, c)
		return
	}

	c.JSON(code, result)
}

// GetCateringClientOrders returns orders of provided client
// @Summary returns orders of provided client
// @Tags caterings orders
// @Produce json
// @Param id path string true "Catering ID"
// @Param clientId path string true "Client ID"
// @Param date query string true "Date query in YYYY-MM-DDT00:00:00Z format"
// @Success 200 {object} response.SummaryOrdersResponse false "Orders for clients"
// @Failure 400 {object} types.Error "Error"
// @Failure 404 {object} types.Error "Not Found"
// @Router /caterings/{id}/clients/{clientId}/orders [get]
func (o Order) GetCateringClientOrders(c *gin.Context) {
	var path types.PathClient
	var query types.DateQuery

	if err := utils.RequestBinderURI(&path, c); err != nil {
		return
	}

	if err := utils.RequestBinderQuery(&query, c); err != nil {
		return
	}

	client := types.CompanyTypesEnum.Catering

	result, code, err := services.NewOrderService().GetCateringOrderService(path, query, client)

	if err != nil {
		utils.CreateError(code, err, c)
		return
	}

	c.JSON(code, result)
}

// ApproveOrders changes status of orders for provided day
// to approved
// @Summary approve user orders
// @Tags clients orders
// @Produce json
// @Param id path string true "Client ID"
// @Param date query string true "Date query in YYYY-MM-DDT00:00:00Z format"
// @Success 200 "Successfully Approved orders"
// @Failure 400 {object} types.Error "Error"
// @Failure 404 {object} types.Error "Not Found"
// @Router /clients/{id}/orders [put]
func (o Order) ApproveOrders(c *gin.Context) {
	var path types.PathID
	var query types.DateQuery

	if err := utils.RequestBinderURI(&path, c); err != nil {
		return
	}

	if err := utils.RequestBinderQuery(&query, c); err != nil {
		return
	}

	code, err := services.NewOrderService().ApproveOrdersService(path, query)
	if err != nil {
		utils.CreateError(code, err, c)
		return
	}

	c.Status(http.StatusOK)
}

// GetOrderStatus returns status of order
// @Summary returns status of order
// @Tags clients orders
// @Produce json
// @Param id path string true "Client ID"
// @Param date query string true "Date query in YYYY-MM-DDT00:00:00Z format"
// @Success 200 {object} response.OrderStatus "order status"
// @Failure 400 {object} types.Error "Error"
// @Failure 404 {object} types.Error "Not Found"
// @Router /clients/{id}/order-status [get]
func (o Order) GetOrderStatus(c *gin.Context) {
	var path types.PathID
	var query types.DateQuery

	if err := utils.RequestBinderURI(&path, c); err != nil {
		return
	}

	if err := utils.RequestBinderQuery(&query, c); err != nil {
		return
	}

	status := services.NewOrderService().GetOrderStatus(path, query)

	c.JSON(http.StatusOK, gin.H{
		"status": status,
	})
}
