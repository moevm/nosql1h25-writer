package get_orders

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/moevm/nosql1h25-writer/backend/internal/api"
	"github.com/moevm/nosql1h25-writer/backend/internal/api/common/decorator"
	"github.com/moevm/nosql1h25-writer/backend/internal/service/orders"
)

type handler struct {
	orderService orders.Service
}

func New(orderService orders.Service) api.Handler {
	return decorator.NewBindAndValidate(&handler{orderService: orderService})
}

type Request struct {
	Offset *int `query:"offset" validate:"gte=0" example:"0"`
	Limit  *int `query:"limit" validate:"gte=1,lte=200" example:"10"`
}

type Response struct {
	Orders []Order `json:"orders"`
	Total  int     `json:"total" example:"250"`
}

type Order struct {
	Title          string  `json:"title"`
	Description    string  `json:"description"`
	CompletionTime int     `json:"completionTime"`
	Cost           int     `json:"cost,omitempty"`
	ClientName     string  `json:"clientName"`
	Rating         float64 `json:"rating"`
}

// @Description	Get a paginated list of orders and total count
// @Summary	Get orders list
// @Tags orders
// @Param offset query int false "Offset" default(0) minimum(0) example(0)
// @Param limit query int false	"Limit" default(10) minimum(1) maximum(200) example(10)
// @Accept json
// @Produce	json
// @Success	200	{object} Response
// @Failure	500	{object} echo.HTTPError
// @Router /orders [get]
func (h *handler) Handle(c echo.Context, in Request) error {
	offset, limit := applyDefaults(in)
	findOut, err := h.orderService.Find(c.Request().Context(), offset, limit)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	orderList := make([]Order, 0, len(findOut.Orders))
	for _, order := range findOut.Orders {
		orderList = append(orderList, Order{
			Title:          order.Title,
			Description:    order.Description,
			CompletionTime: order.CompletionTime,
			Cost:           order.Cost,
			ClientName:     order.ClientName,
			Rating:         order.Rating,
		})
	}
	return c.JSON(http.StatusOK, Response{Orders: orderList, Total: len(orderList)})
}

func applyDefaults(in Request) (offset int, limit int) {
	if in.Offset != nil {
		offset = *in.Offset
	} else {
		offset = 0
	}
	if in.Limit != nil {
		limit = *in.Limit
	} else {
		limit = 10
	}
	return
}
