package get_users_id_orders

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/moevm/nosql1h25-writer/backend/internal/api"
	"github.com/moevm/nosql1h25-writer/backend/internal/api/common/decorator"
	"github.com/moevm/nosql1h25-writer/backend/internal/api/common/mw"
	"github.com/moevm/nosql1h25-writer/backend/internal/entity"
	"github.com/moevm/nosql1h25-writer/backend/internal/service/users"
)

type handler struct {
	users users.Service
}

type Request struct {
	ID primitive.ObjectID `param:"id" validate:"required" example:"507f1f77bcf86cd799439011"`
}

func New(users users.Service) api.Handler {
	return decorator.NewBindAndValidate(&handler{users: users})
}

// @Description	Получить список заказов пользователя
// @Summary		Получить список заказов пользователя
// @Tags			Users
// @Security		JWT
// @Param			id	path	string	true	"ID пользователя"	example("507f1f77bcf86cd799439011")
// @Produce		json
// @Success		200	{array}		entity.OrderExt
// @Failure		400	{object}	echo.HTTPError
// @Failure		403	{object}	echo.HTTPError
// @Failure		404	{object}	echo.HTTPError
// @Failure		500	{object}	echo.HTTPError
// @Router			/users/{id}/orders [get]
func (h *handler) Handle(c echo.Context, in Request) error {
	requesterID := c.Get(mw.UserIDKey).(primitive.ObjectID) //nolint:forcetypeassert
	role := c.Get(mw.SystemRoleKey).(entity.SystemRoleType) //nolint:forcetypeassert
	isAdmin := role == entity.SystemRoleTypeAdmin

	orders, err := h.users.FindOrdersByUserID(c.Request().Context(), requesterID, in.ID, isAdmin)
	if err != nil {
		switch {
		case errors.Is(err, users.ErrForbidden):
			return echo.NewHTTPError(http.StatusForbidden, "access denied")
		case errors.Is(err, users.ErrUserNotFound):
			return echo.NewHTTPError(http.StatusNotFound, "user not found")
		default:
			return echo.NewHTTPError(http.StatusInternalServerError, "internal server error")
		}
	}

	return c.JSON(http.StatusOK, orders)
}
