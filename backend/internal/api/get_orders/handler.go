package get_orders

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/moevm/nosql1h25-writer/backend/internal/api"
	"github.com/moevm/nosql1h25-writer/backend/internal/api/common/decorator"
	orders_repo "github.com/moevm/nosql1h25-writer/backend/internal/repo/orders"
	"github.com/moevm/nosql1h25-writer/backend/internal/service/orders"
)

type handler struct {
	orderService orders.Service
}

func New(orderService orders.Service) api.Handler {
	return decorator.NewBindAndValidate(&handler{orderService: orderService})
}

type Request struct {
	Offset int64 `query:"offset" validate:"gte=0" example:"0"`
	Limit  int64 `query:"limit" validate:"gte=1,lte=100" example:"10"`
}

type Response struct {
	Orders []orders_repo.FindOrdersOut `json:"orders"`
	Total  int64                       `json:"total" example:"250"`
}

func (h *handler) ApplyDefaults(req *Request) {
	if req.Limit == 0 {
		req.Limit = 10
	}
	if req.Offset < 0 {
		req.Offset = 0
	}
}

// @Description	Get a paginated list of orders and total count
// @Summary	Get orders list
// @Tags orders
// @Param offset query int false "Offset" default(0) minimum(0) example(0)
// @Param limit query int false	"Limit" default(10) minimum(1) maximum(100) example(10)
// @Accept json
// @Produce	json
// @Success	200	{object} Response
// @Failure 400 {object} echo.HTTPError
// @Failure	500	{object} echo.HTTPError
// @Router /orders [get]
func (h *handler) Handle(c echo.Context, in Request) error {
	ctx := c.Request().Context()
	h.ApplyDefaults(&in)

	result, total, err := h.orderService.FindOrders(ctx, in.Offset, in.Limit)
	if err != nil {
		if errors.Is(err, orders.ErrOrdersNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, err.Error())
		} else if errors.Is(err, orders.ErrInvalidPagination) {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	orders := make([]orders_repo.FindOrdersOut, len(result))
	for i, o := range result {
		orders[i] = orders_repo.FindOrdersOut{
			ID:           o.ID,
			Title:        o.Title,
			Description:  o.Description,
			Budget:       o.Budget,
			Active:       o.Active,
			CreatedAt:    o.CreatedAt,
			ClientID:     o.ClientID,
			FreelancerID: o.FreelancerID,
		}
	}

	return c.JSON(http.StatusOK, Response{
		Orders: orders,
		Total:  total,
	})
}
