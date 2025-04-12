package post_balance_deposit

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/moevm/nosql1h25-writer/backend/internal/api"
	"github.com/moevm/nosql1h25-writer/backend/internal/api/common/decorator"
	"github.com/moevm/nosql1h25-writer/backend/internal/api/common/mw"
	"github.com/moevm/nosql1h25-writer/backend/internal/service/users"
)

type handler struct {
	usersService users.Service
}

func New(usersService users.Service) api.Handler {
	return decorator.NewBindAndValidate(&handler{usersService: usersService})
}

type Request struct {
	Amount int `json:"amount" validate:"required,gt=0" minimum:"1" example:"100"`
}

type Response struct {
	NewBalance int `json:"newBalance" example:"777"`
}

// Handle - Deposit funds to user's balance
//
//	@Summary		Deposit funds
//	@Description	Add specified amount to authenticated user's balance
//	@Tags			balance
//	@Security		JWT
//	@Param			request	body	Request	true	"Deposit amount"
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	Response
//	@Failure		400	{object}	echo.HTTPError
//	@Failure		401	{object}	echo.HTTPError
//	@Failure		500	{object}	echo.HTTPError
//	@Router			/balance/deposit [post]
func (h *handler) Handle(c echo.Context, in Request) error {
	userID := c.Get(mw.UserIDKey).(primitive.ObjectID) //nolint:forcetypeassert

	newBalance, err := h.usersService.UpdateBalance(c.Request().Context(), userID, users.OperationTypeDeposit, in.Amount)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, Response{NewBalance: newBalance})
}
