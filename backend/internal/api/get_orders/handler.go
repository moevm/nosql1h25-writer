package get_orders

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/moevm/nosql1h25-writer/backend/internal/api"
	"github.com/sv-tools/mongoifc"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type handler struct {
	orders mongoifc.Collection
}

func New(orders mongoifc.Collection) api.Handler {
	return &handler{orders: orders}
}

// @Description	Получить список заказов с пагинацией и общее кол-во заказов
// @Summary	Получить список заказов
// @Tags orders
// @Param offset query int false "Offset" default(0) minimum(0) example(0)
// @Param limit query int false	"Limit" default(10) minimum(1) maximum(100) example(10)
// @Produce	json
// @Success	200	{object} map[string]interface{}
// @Failure	500	{object} echo.HTTPError
// @Router /orders [get]
func (h *handler) Handle(c echo.Context) error {
	offset, _ := strconv.Atoi(c.QueryParam("offset"))
	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	if limit == 0 {
		limit = 10
	}

	count, err := h.orders.CountDocuments(c.Request().Context(), bson.M{})
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	cursor, err := h.orders.Find(c.Request().Context(), bson.M{}, &options.FindOptions{
		Skip:  int64Ptr(int64(offset)),
		Limit: int64Ptr(int64(limit)),
	})
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	defer cursor.Close(c.Request().Context())

	var results []map[string]interface{}
	if err := cursor.All(c.Request().Context(), &results); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"total":  count,
		"orders": results,
	})
}

func int64Ptr(i int64) *int64 {
	return &i
}
