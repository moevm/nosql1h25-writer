package get_users_id_responses

import (
	"net/http"
	"time"

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

type ResponseOrder struct {
	OrderID        primitive.ObjectID `json:"orderId"`
	Title          string             `json:"title"`
	CompletionTime int64              `json:"completionTime"`
	Cost           int                `json:"cost,omitempty"`
	Status         entity.StatusType  `json:"status"`
	CoverLetter    string             `json:"coverLetter"`
	CreatedAt      time.Time          `json:"createdAt"`
}

type Response struct {
	Responses []ResponseOrder `json:"responses"`
}

func New(users users.Service) api.Handler {
	return decorator.NewBindAndValidate(&handler{users: users})
}

// Handle - find orders by response userID
//
// @Description	Получить список заказов, на которые откликался пользователь
// @Summary		Получить список заказов, на которые откликался пользователь
// @Tags			Users
// @Security		JWT
// @Param			id	path	string	true	"ID пользователя"	example("507f1f77bcf86cd799439011")
// @Produce		json
// @Success		200	{object}	Response
// @Failure		400	{object}	echo.HTTPError
// @Failure		403	{object}	echo.HTTPError
// @Failure		404	{object}	echo.HTTPError
// @Failure		500	{object}	echo.HTTPError
// @Router			/users/{id}/responses [get]
func (h *handler) Handle(c echo.Context, in Request) error {
	requesterID := c.Get(mw.UserIDKey).(primitive.ObjectID) //nolint:forcetypeassert
	role := c.Get(mw.SystemRoleKey).(entity.SystemRoleType) //nolint:forcetypeassert

	if requesterID != in.ID && role != entity.SystemRoleTypeAdmin {
		return echo.NewHTTPError(http.StatusForbidden, "access denied")
	}

	ordersExt, err := h.users.FindOrdersByResponseUserID(c.Request().Context(), in.ID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	responses := make([]ResponseOrder, 0, len(ordersExt))
	for _, orderExt := range ordersExt {
		var userResponse entity.Response
		for _, response := range orderExt.Responses {
			if response.FreelancerID == in.ID {
				userResponse = response
				break
			}
		}

		responses = append(responses, ResponseOrder{
			OrderID:        orderExt.ID,
			Title:          orderExt.Title,
			CompletionTime: orderExt.CompletionTime,
			Cost:           orderExt.Cost,
			Status:         orderExt.Statuses[len(orderExt.Statuses)-1].Type,
			CoverLetter:    userResponse.CoverLetter,
			CreatedAt:      userResponse.CreatedAt,
		})
	}

	return c.JSON(http.StatusOK, Response{Responses: responses})
}
