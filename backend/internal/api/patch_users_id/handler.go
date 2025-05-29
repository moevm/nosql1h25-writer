package patch_users_id

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
	usersService users.Service
}

func New(usersService users.Service) api.Handler {
	return decorator.NewBindAndValidate(&handler{usersService: usersService})
}

type Request struct {
	ID                    primitive.ObjectID `param:"id" validate:"required"`
	DisplayName           *string            `json:"displayName" validate:"omitempty,min=3,max=64" example:"username"`
	FreelancerDescription *string            `json:"freelancerDescription,omitempty"`
	ClientDescription     *string            `json:"clientDescription,omitempty"`
}

type Response struct {
	Updated bool `json:"updated" example:"true"`
}

// Handle - Update user
//
//	@Summary		Update user
//	@Description	Partially update user fields. Admin can update any user; regular user can update only their own profile.
//	@Tags			Users
//	@Security		JWT
//	@Param			id		path	string	true	"User ID"
//	@Param			request	body	Request	true	"Fields to update"
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	Response
//	@Failure		400	{object}	echo.HTTPError
//	@Failure		401	{object}	echo.HTTPError
//	@Failure		403	{object}	echo.HTTPError
//	@Failure		404	{object}	echo.HTTPError
//	@Failure		500	{object}	echo.HTTPError
//	@Router			/users/{id} [patch]
func (h *handler) Handle(c echo.Context, in Request) error {
	RequesterRole := c.Get(mw.SystemRoleKey).(entity.SystemRoleType) //nolint:forcetypeassert
	RequesterID := c.Get(mw.UserIDKey).(primitive.ObjectID)          //nolint:forcetypeassert

	// Проверка: если не передано ни одного изменяемого поля
	if in.DisplayName == nil && in.FreelancerDescription == nil &&
		in.ClientDescription == nil {
		return echo.NewHTTPError(http.StatusBadRequest, "no fields to update")
	}

	update := users.UpdateInput{
		RequesterID:           RequesterID,
		RequesterRole:         RequesterRole,
		UserID:                in.ID,
		DisplayName:           in.DisplayName,
		FreelancerDescription: in.FreelancerDescription,
		ClientDescription:     in.ClientDescription,
	}

	err := h.usersService.UpdateProfile(c.Request().Context(), update)
	if err != nil {
		if errors.Is(err, users.ErrUserNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, "user not found")
		}
		if errors.Is(err, users.ErrForbidden) {
			return echo.NewHTTPError(http.StatusForbidden, "access denied")
		}

		return echo.NewHTTPError(http.StatusInternalServerError, "internal error")
	}

	return c.JSON(http.StatusOK, Response{Updated: true})
}
