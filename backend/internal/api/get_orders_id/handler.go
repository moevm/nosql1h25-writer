package get_orders_id

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/moevm/nosql1h25-writer/backend/internal/api"
	"github.com/moevm/nosql1h25-writer/backend/internal/api/common/decorator"
	"github.com/moevm/nosql1h25-writer/backend/internal/api/common/mw"
	"github.com/moevm/nosql1h25-writer/backend/internal/entity"
	"github.com/moevm/nosql1h25-writer/backend/internal/service/orders"
	"github.com/moevm/nosql1h25-writer/backend/internal/service/users"
)

type handler struct {
	orderService orders.Service
	usersService users.Service
}

func New(orderService orders.Service, usersService users.Service) api.Handler {
	return decorator.NewBindAndValidate(&handler{
		orderService: orderService,
		usersService: usersService,
	})
}

type Request struct {
	ID primitive.ObjectID `param:"id" validate:"required"`
}

type Order struct {
	ID             primitive.ObjectID `json:"id" validate:"required" example:"582ebf010936ac3ba5cd00e4"`
	ClientName     string             `json:"clientName" validate:"required" example:"John Doe"`
	ClientRating   float64            `json:"clientRating" validate:"required" example:"4.8"`
	ClientID       primitive.ObjectID `json:"clientId" validate:"required" example:"582ebf010936ac3ba5cd00e4"`
	FreelancerID   primitive.ObjectID `json:"freelancerId,omitempty" validate:"omitempty" example:"582ebf010936ac3ba5cd00e4"`
	Status         entity.StatusType  `json:"status" validate:"required" example:"beginning"`
	Title          string             `json:"title" validate:"required" example:"Write something for me"`
	Description    string             `json:"description" validate:"required" example:"Write something for me but more words"`
	CompletionTime int64              `json:"completionTime" validate:"required" example:"3600000000000"`
	Cost           int                `json:"cost,omitempty" validate:"omitempty" example:"500"`
	CreatedAt      time.Time          `json:"createdAt" validate:"required" example:"2020-01-01T00:00:00Z"`
	UpdatedAt      time.Time          `json:"updatedAt" validate:"required" example:"2020-01-01T00:00:00Z"`
}

type Response struct {
	Order             Order `json:"order" validate:"required"`
	HasActiveResponse bool  `json:"hasActiveResponse" validate:"required"`
	IsClient          bool  `json:"isClient" validate:"required"`
	IsFreelancer      bool  `json:"isFreelancer" validate:"required"`
}

// Handle - Get Order by ID
//
//	@Description	Return order by ID
//	@Summary		Get info about order and several related things
//	@Tags			orders
//	@Security		JWT
//	@Param			id	path	string	true	"Order ID"
//	@Produce		json
//	@Success		200	{object}	Response
//	@Failure		401	{object}	echo.HTTPError
//	@Failure		400	{object}	echo.HTTPError
//	@Failure		404	{object}	echo.HTTPError
//	@Failure		500	{object}	echo.HTTPError
//	@Router			/orders/{id} [get]
func (h *handler) Handle(c echo.Context, in Request) error {
	order, err := h.orderService.GetByIDExt(c.Request().Context(), in.ID)
	if err != nil {
		if errors.Is(err, orders.ErrOrderNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, err.Error())
		}

		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	client, err := h.usersService.GetByIDExt(c.Request().Context(), order.ClientID)
	if err != nil {
		if errors.Is(err, users.ErrUserNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, fmt.Sprintf("client fetching: %v", err))
		}

		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("client fetching: %v", err))
	}

	userID := c.Get(mw.UserIDKey).(primitive.ObjectID) //nolint:forcetypeassert
	res := Response{
		Order: Order{
			ID:             order.ID,
			ClientName:     client.DisplayName,
			ClientRating:   client.Client.Rating,
			ClientID:       order.ClientID,
			FreelancerID:   order.FreelancerID,
			Status:         order.Statuses[len(order.Statuses)-1].Type,
			Title:          order.Title,
			Description:    order.Description,
			CompletionTime: order.CompletionTime,
			Cost:           order.Cost,
			CreatedAt:      order.CreatedAt,
			UpdatedAt:      order.UpdatedAt,
		},
		IsClient:     userID == order.ClientID,
		IsFreelancer: userID == order.FreelancerID,
	}

	for _, response := range order.Responses {
		if response.Active && response.FreelancerID == userID {
			res.HasActiveResponse = true
			break
		}
	}

	return c.JSON(http.StatusOK, res)
}
