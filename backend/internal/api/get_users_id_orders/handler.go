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

type Order struct {
	ID             primitive.ObjectID `json:"id"`
	ClientID       primitive.ObjectID `json:"clientId"`
	Title          string             `json:"title"`
	Description    string             `json:"description"`
	CompletionTime int64              `json:"completionTime"`
	Status         entity.StatusType  `json:"status"`
	TotalResponses int                `json:"totalResponses"`
	Cost           int                `json:"cost"`
	FreelancerID   primitive.ObjectID `json:"freelancerId"`
	CreatedAt      string             `json:"createdAt"`
	UpdatedAt      string             `json:"updatedAt"`
}

type Response struct {
	Orders []Order `json:"orders"`
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
// @Success		200	{object}	Response
// @Failure		400	{object}	echo.HTTPError
// @Failure		403	{object}	echo.HTTPError
// @Failure		404	{object}	echo.HTTPError
// @Failure		500	{object}	echo.HTTPError
// @Router			/users/{id}/orders [get]
func (h *handler) Handle(c echo.Context, in Request) error {
	requesterID := c.Get(mw.UserIDKey).(primitive.ObjectID) //nolint:forcetypeassert
	role := c.Get(mw.SystemRoleKey).(entity.SystemRoleType) //nolint:forcetypeassert

	if requesterID != in.ID && role != entity.SystemRoleTypeAdmin {
		return echo.NewHTTPError(http.StatusForbidden, "access denied")
	}

	ordersExt, err := h.users.FindOrdersByUserID(c.Request().Context(), requesterID, in.ID)
	if err != nil {
		if errors.Is(err, users.ErrUserNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, "user not found")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "internal server error")
	}

	orders := make([]Order, 0, len(ordersExt))
	for _, orderExt := range ordersExt {
		status := entity.StatusTypeBeginning
		if len(orderExt.Statuses) > 0 {
			status = orderExt.Statuses[len(orderExt.Statuses)-1].Type
		}

		totalResponses := 0
		for _, response := range orderExt.Responses {
			if response.Active {
				totalResponses++
			}
		}

		orders = append(orders, Order{
			ID:             orderExt.ID,
			ClientID:       orderExt.ClientID,
			Title:          orderExt.Title,
			Description:    orderExt.Description,
			CompletionTime: orderExt.CompletionTime,
			Status:         status,
			TotalResponses: totalResponses,
			Cost:           orderExt.Cost,
			FreelancerID:   orderExt.FreelancerID,
			CreatedAt:      orderExt.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt:      orderExt.UpdatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	return c.JSON(http.StatusOK, Response{Orders: orders})
}
