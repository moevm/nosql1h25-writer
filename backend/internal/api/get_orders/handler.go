package get_orders

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"

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
	Offset  *int    `query:"offset" validate:"gte=0" example:"0"`
	Limit   *int    `query:"limit" validate:"gte=1,lte=200" example:"10"`
	MinCost *int    `query:"minCost" validate:"omitempty,gte=0" example:"100"`
	MaxCost *int    `query:"maxCost" validate:"omitempty,gte=0" example:"1000"`
	MinTime *int64  `query:"minTime" validate:"omitempty,gte=3600000000000" example:"3600000000000"`
	MaxTime *int64  `query:"maxTime" validate:"omitempty,gte=3600000000000" example:"3600000000000"`
	Search  *string `query:"search" validate:"omitempty" example:"Написать сценарий"`
	SortBy  *string `query:"sortBy" validate:"omitempty,oneof=newest oldest cost_asc cost_desc time_asc time_desc" example:"newest"`
}

type Order struct {
	ID             primitive.ObjectID `json:"id"`
	Title          string             `json:"title"`
	Description    string             `json:"description"`
	CompletionTime int64              `json:"completionTime"`
	Cost           int                `json:"cost,omitempty"`
}

type Response struct {
	Orders []Order `json:"orders"`
	Total  int     `json:"total" example:"250"`
}

// Handle - Get Orders
//
//	@Description	Get a paginated list of orders and total count
//	@Summary		Get orders list
//	@Tags			orders
//	@Security		JWT
//	@Param			offset	query	int		false	"Offset"											default(0)	minimum(0)	example(0)
//	@Param			limit	query	int		false	"Limit"												default(10)	minimum(1)	maximum(200)	example(10)
//	@Param			sortBy	query	string	false	"Sort field: cost_asc, cost_desc, newest, oldest"	Enums(cost_asc,cost_desc,newest,oldest)
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	Response
//	@Failure		500	{object}	echo.HTTPError
//	@Router			/orders [get]
func (h *handler) Handle(c echo.Context, in Request) error {
	offset, limit := h.applyDefaults(in)

	out, err := h.orderService.Find(c.Request().Context(), orders.FindIn{
		Limit:   limit,
		Offset:  offset,
		MinCost: in.MinCost,
		MaxCost: in.MaxCost,
		MinTime: in.MinTime,
		MaxTime: in.MaxTime,
		Search:  in.Search,
		SortBy:  in.SortBy,
	})
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	orders := make([]Order, 0, len(out.Orders))
	for _, order := range out.Orders {
		orders = append(orders, Order{
			ID:             order.ID,
			Title:          order.Title,
			Description:    order.Description,
			CompletionTime: order.CompletionTime,
			Cost:           order.Cost,
		})
	}

	return c.JSON(http.StatusOK, Response{Orders: orders, Total: out.Total})
}

func (h *handler) applyDefaults(in Request) (offset int, limit int) {
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
