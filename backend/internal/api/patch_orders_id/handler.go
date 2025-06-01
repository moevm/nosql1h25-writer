package patch_orders_id

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/moevm/nosql1h25-writer/backend/internal/api"
	"github.com/moevm/nosql1h25-writer/backend/internal/api/common/decorator"
	"github.com/moevm/nosql1h25-writer/backend/internal/api/common/mw"
	"github.com/moevm/nosql1h25-writer/backend/internal/entity"
	"github.com/moevm/nosql1h25-writer/backend/internal/service/orders"
)

type handler struct {
	ordersService orders.Service
}

func New(ordersService orders.Service) api.Handler {
	return decorator.NewBindAndValidate(&handler{ordersService: ordersService})
}

type Request struct {
	ID             primitive.ObjectID  `param:"id" validate:"required"`
	Title          *string             `json:"title,omitempty" validate:"omitempty,min=4,max=256" example:"New title"`
	Description    *string             `json:"description,omitempty" validate:"omitempty,min=16,max=2048" example:"New Order Description"`
	CompletionTime *int64              `json:"completionTime,omitempty" validate:"omitempty,gte=3600000000000" example:"3600000000000"`
	Cost           *int                `json:"cost,omitempty" validate:"omitempty,min=0" example:"5000"`
	Status         *entity.StatusType  `json:"status,omitempty" validate:"omitempty,oneof=beginning negotiation budgeting work reviews finished dispute" example:"finished"`
	FreelancerID   *primitive.ObjectID `json:"freelancerId,omitempty"`
}

// Handle - Update order
//
//	@Summary		Update order
//	@Description	Only updates fields present in the request. Admin can update any order. User can update only their own open orders.
//	@Tags			orders
//	@Security		JWT
//	@Param			id		path	string	true	"Order ID"
//	@Param			request	body	Request	true	"Fields to update"
//	@Accept			json
//	@Produce		json
//	@Success		200
//	@Failure		400	{object}	echo.HTTPError
//	@Failure		401	{object}	echo.HTTPError
//	@Failure		403	{object}	echo.HTTPError
//	@Failure		404	{object}	echo.HTTPError
//	@Failure		500	{object}	echo.HTTPError
//	@Router			/orders/{id} [patch]
func (h *handler) Handle(c echo.Context, in Request) error {
	role := c.Get(mw.SystemRoleKey).(entity.SystemRoleType) //nolint:forcetypeassert
	userID := c.Get(mw.UserIDKey).(primitive.ObjectID)      //nolint:forcetypeassert

	// Проверка: если не передано ни одного изменяемого поля
	if in.Title == nil &&
		in.Description == nil &&
		in.CompletionTime == nil &&
		in.Cost == nil &&
		in.Status == nil &&
		in.FreelancerID == nil {
		return echo.NewHTTPError(http.StatusBadRequest, "no fields to update")
	}

	// Если не админ, получаем заказ и проверяем является ли юзер заказчиком, указанного заказа
	if role != entity.SystemRoleTypeAdmin {
		order, err := h.ordersService.GetByID(c.Request().Context(), in.ID)
		if err != nil {
			if errors.Is(err, orders.ErrOrderNotFound) {
				return echo.NewHTTPError(http.StatusNotFound, err.Error())
			}

			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		if order.ClientID != userID {
			return echo.NewHTTPError(http.StatusForbidden, "access denied")
		}
	}

	err := h.ordersService.Update(c.Request().Context(), orders.UpdateIn{
		OrderID:        in.ID,
		Title:          in.Title,
		Description:    in.Description,
		CompletionTime: in.CompletionTime,
		Cost:           in.Cost,
		Status:         in.Status,
		FreelancerID:   in.FreelancerID,
	})
	if err != nil {
		if errors.Is(err, orders.ErrOrderNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, err.Error())
		}

		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.NoContent(http.StatusOK)
}
