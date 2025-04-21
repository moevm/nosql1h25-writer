package get_orders_id

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/moevm/nosql1h25-writer/backend/internal/api"
	"github.com/moevm/nosql1h25-writer/backend/internal/api/common/decorator"
	"github.com/moevm/nosql1h25-writer/backend/internal/service/orders"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type handler struct {
	orderService orders.Service
}

func New(orderService orders.Service) api.Handler {
	return decorator.NewBindAndValidate(&handler{orderService: orderService})
}

type Request struct {
	ID primitive.ObjectID `param:"id" validate:"required"`
}

type Response struct {
	Title          string  `json:"title"`
	Description    string  `json:"description"`
	CompletionTime int     `json:"completionTime"`
	Cost           int     `json:"cost,omitempty"`
	ClientName     string  `json:"clientName"`
	Rating         float64 `json:"rating"`
}

// Handle - Get Ordes by ID
//
//	@Description	Return order by MongoDB ObjectID
//	@Summary		Get info about order
//	@Tags			orders
//	@Security		JWT
//	@Param			id	path	string	true	"Order ID"
//	@Produce		json
//	@Success		200	{object}	map[string]interface{}
//	@Failure		400	{object}	echo.HTTPError	"Incorrect ID"
//	@Failure		404	{object}	echo.HTTPError	"Order not found"
//	@Failure		500	{object}	echo.HTTPError
//	@Router			/orders/{id} [get]
func (h *handler) Handle(c echo.Context, in Request) error {
	order, err := h.orderService.GetByID(c.Request().Context(), in.ID)
	if err != nil {
		if errors.Is(err, orders.ErrOrderNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, err.Error())
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, Response{
		Title:          order.Title,
		Description:    order.Description,
		CompletionTime: order.CompletionTime,
		Cost:           order.Cost,
		ClientName:     order.ClientName,
		Rating:         order.Rating,
	})
}
