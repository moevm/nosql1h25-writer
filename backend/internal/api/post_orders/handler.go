package post_orders

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
	Title          string `json:"title" validate:"required,min=3,max=32" example:"Сценарий"`
	Description    string `json:"description" validate:"required,min=16,max=8192" example:"Написать сценарий вот такой и такой"`
	CompletionTime int64  `json:"completionTime" validate:"required,gte=3600000000000" example:"3600000000000"`
	Cost           int    `json:"cost" validate:"gte=0" example:"500"`
}

type Response struct {
	ID primitive.ObjectID `json:"id" validate:"required" example:"522bb79455449d881b004d27"`
}

// Handle - Create order handler
//
//	@Summary		Create order
//	@Description	Create order on behalf of authenticated user
//	@Tags			orders
//	@Security		JWT
//	@Param			request	body	Request	true "order parameters"
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	Response
//	@Failure		400	{object}	echo.HTTPError
//	@Failure		401	{object}	echo.HTTPError
//	@Failure		500	{object}	echo.HTTPError
//	@Router			/orders [post]
func (h *handler) Handle(c echo.Context, in Request) error {
	id, err := h.ordersService.Create(c.Request().Context(), orders.CreateIn{ //nolint:forcetypeassert
		ClientID:       c.Get(mw.UserIDKey).(primitive.ObjectID),
		Title:          in.Title,
		Description:    in.Description,
		CompletionTime: in.CompletionTime,
		Cost:           in.Cost,
	})

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, Response{ID: id})
}
