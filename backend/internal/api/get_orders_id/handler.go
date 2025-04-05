package get_orders_id

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/moevm/nosql1h25-writer/backend/internal/api"
	"github.com/sv-tools/mongoifc"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type handler struct {
	orders mongoifc.Collection
}

func New(orders mongoifc.Collection) api.Handler {
	return &handler{orders: orders}
}

// @Description	Возвращает один заказ по его MongoDB ObjectID
// @Summary Получить информацию о заказе
// @Tags orders
// @Param id path string true "Order ID"
// @Produce	json
// @Success	200	{object} map[string]interface{}
// @Failure	400	{object} echo.HTTPError "Неверный формат ID"
// @Failure	404	{object} echo.HTTPError "Заказ не найден"
// @Failure	500	{object} echo.HTTPError
// @Router /orders/{id} [get]
func (h *handler) Handle(c echo.Context) error {
	idStr := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid ID format")
	}

	var order map[string]interface{}
	err = h.orders.FindOne(c.Request().Context(), bson.M{"_id": id}).Decode(&order)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Order not found")
	}

	return c.JSON(http.StatusOK, order)
}
