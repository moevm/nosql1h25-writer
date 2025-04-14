package get_orders_id

// import (
// 	"errors"
// 	"net/http"

// 	"github.com/labstack/echo/v4"

// 	"github.com/moevm/nosql1h25-writer/backend/internal/api"
// 	"github.com/moevm/nosql1h25-writer/backend/internal/service/orders"
// )

// type handler struct {
// 	orderService orders.Service
// }

// func New(orderService orders.Service) api.Handler {
// 	return &handler{orderService: orderService}
// }

// // @Description	Return order by MongoDB ObjectID
// // @Summary Get info about order
// // @Tags orders
// // @Param id path string true "Order ID"
// // @Produce	json
// // @Success	200	{object} map[string]interface{}
// // @Failure	400	{object} echo.HTTPError "Incorrect ID"
// // @Failure	404	{object} echo.HTTPError "Order not found"
// // @Failure	500	{object} echo.HTTPError
// // @Router /orders/{id} [get]
// func (h *handler) Handle(c echo.Context) error {
// 	id := c.Param("id")
// 	order, err := h.orderService.GetByID(c.Request().Context(), id)
// 	if err != nil {
// 		if errors.Is(err, orders.ErrOrdersNotFound) {
// 			return echo.NewHTTPError(http.StatusNotFound, err.Error)
// 		}
// 		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
// 	}
// 	return c.JSON(http.StatusOK, order)
// }
