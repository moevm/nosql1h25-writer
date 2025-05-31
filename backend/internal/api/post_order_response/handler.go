package post_order_response

import (
	"errors"
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
	OrderID     primitive.ObjectID `param:"id" validate:"required" example:"683b2dc10949bd1e64266ed0"`
	CoverLetter string             `json:"coverLetter" validate:"required" example:"Я заинтересован в вашем проекте и имею релевантный опыт в этой области. Готов обсудить детали и начать работу."`
}

type Response struct {
	ID primitive.ObjectID `json:"id" validate:"required" example:"683b2dc10949bd1e64266ed0"`
}

// Handle - Response to order handler
//
//	@Summary		Response to order
//	@Description	Create response to existing order
//	@Tags			orders
//	@Security		JWT
//	@Param			id		path	string	true	"Order ID"		example(683b2dc10949bd1e64266ed0)
//	@Param			request	body	Request	true	"Response data"	example({"coverLetter": "Я заинтересован в вашем проекте и имею релевантный опыт в этой области. Готов обсудить детали и начать работу."})
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	Response
//	@Failure		400	{object}	echo.HTTPError
//	@Failure		401	{object}	echo.HTTPError
//	@Failure		404	{object}	echo.HTTPError
//	@Failure		500	{object}	echo.HTTPError
//	@Router			/orders/{id}/response [post]
func (h *handler) Handle(c echo.Context, in Request) error {
	userID, ok := c.Get(mw.UserIDKey).(primitive.ObjectID)
	if !ok {
		return echo.NewHTTPError(http.StatusInternalServerError, "invalid user id")
	}

	err := h.ordersService.CreateResponse(c.Request().Context(), in.OrderID, userID, in.CoverLetter)
	if err != nil {
		switch {
		case errors.Is(err, orders.ErrOrderNotFound):
			return echo.NewHTTPError(http.StatusNotFound, "order not found")
		case errors.Is(err, orders.ErrCannotResponse):
			return echo.NewHTTPError(http.StatusBadRequest, "cannot response to this order")
		default:
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
	}

	return c.JSON(http.StatusOK, Response{ID: in.OrderID})
}
