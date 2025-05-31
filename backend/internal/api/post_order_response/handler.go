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
	"github.com/moevm/nosql1h25-writer/backend/internal/service/users"
)

type handler struct {
	ordersService orders.Service
}

func New(ordersService orders.Service) api.Handler {
	return decorator.NewBindAndValidate(&handler{ordersService: ordersService})
}

type Request struct {
	OrderID     primitive.ObjectID `param:"id" validate:"required" example:"683b2dc10949bd1e64266ed0"`
	CoverLetter string             `json:"coverLetter" validate:"required,min=16,max=512" example:"Я заинтересован в вашем проекте и имею релевантный опыт в этой области. Готов обсудить детали и начать работу."`
}

// Handle - Response to order handler
//
//	@Summary		Response to order
//	@Description	Create response to existing order
//	@Tags			orders
//	@Security		JWT
//	@Param			id		path	string	true	"Order ID"
//	@Param			request	body	Request	true	"Response data"
//	@Accept			json
//	@Failure		400	{object}	echo.HTTPError
//	@Failure		401	{object}	echo.HTTPError
//	@Failure		404	{object}	echo.HTTPError
//	@Failure		500	{object}	echo.HTTPError
//	@Router			/orders/{id}/response [post]
func (h *handler) Handle(c echo.Context, in Request) error {
	userID := c.Get(mw.UserIDKey).(primitive.ObjectID) //nolint:forcetypeassert

	err := h.ordersService.CreateResponse(c.Request().Context(), in.OrderID, userID, in.CoverLetter)
	if err != nil {
		switch {
		case errors.Is(err, orders.ErrOrderNotFound) || errors.Is(err, users.ErrUserNotFound):
			return echo.NewHTTPError(http.StatusNotFound, err.Error())
		case errors.Is(err, orders.ErrCannotResponse):
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		default:
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
	}

	return c.NoContent(http.StatusOK)
}
