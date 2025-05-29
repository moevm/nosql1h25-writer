package post_order_response

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/moevm/nosql1h25-writer/backend/internal/api"
	"github.com/moevm/nosql1h25-writer/backend/internal/api/common/decorator"
	"github.com/moevm/nosql1h25-writer/backend/internal/api/common/mw"
	"github.com/moevm/nosql1h25-writer/backend/internal/service/orders"
)

type handler struct {
	ordersService orders.Service
}

func New(ordersService orders.Service) api.Handler {
	return decorator.NewBindAndValidate(&handler{ordersService: ordersService})
}

type Request struct {
	OrderID primitive.ObjectID `param:"id" validate:"required" example:"522bb79455449d881b004d27"`
}

type Response struct {
	ID primitive.ObjectID `json:"id" validate:"required" example:"522bb79455449d881b004d27"`
}

// Handle - Response to order handler
//
//	@Summary		Response to order
//	@Description	Create response to existing order
//	@Tags			orders
//	@Security		JWT
//	@Param			id	path	string	true	"Order ID"	example(6838b94ef4b02cca187b2ec2)
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	Response
//	@Failure		400	{object}	echo.HTTPError
//	@Failure		401	{object}	echo.HTTPError
//	@Failure		500	{object}	echo.HTTPError
//	@Router			/orders/{id}/response [post]
func (h *handler) Handle(c echo.Context, in Request) error {
	id, err := h.ordersService.Response(c.Request().Context(), in.OrderID, c.Get(mw.UserIDKey).(primitive.ObjectID))
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, Response{ID: id})
}
