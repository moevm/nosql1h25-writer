package post_auth_logout

import (
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

	"github.com/moevm/nosql1h25-writer/backend/internal/api"
	"github.com/moevm/nosql1h25-writer/backend/internal/api/common/decorator"
	"github.com/moevm/nosql1h25-writer/backend/internal/api/common/refresh"
	"github.com/moevm/nosql1h25-writer/backend/internal/service/auth"
)

type handler struct {
	authService auth.Service
}

func New(authService auth.Service) api.Handler {
	return decorator.NewBindAndValidate(&handler{authService: authService})
}

type Request struct {
	RefreshToken *uuid.UUID `json:"refreshToken" example:"0e8f711e-b713-4869-b528-059a74311482"`
}

// Handle - Logout handler
//
// @Summary		Logout
// @Description	Remove `refreshSession` attached to `refreshToken`. `refreshToken` can be passed in cookie
// @Tags			auth
// @Param			refreshToken	body	Request	false	"active refresh token in UUID RFC4122 format"
// @Accept			json
// @Produce		json
// @Success		200
// @Failure		400	{object}	echo.HTTPError
// @Failure		401	{object}	echo.HTTPError
// @Failure		500	{object}	echo.HTTPError
// @Router			/auth/logout [post]
func (h *handler) Handle(c echo.Context, in Request) error {
	refreshToken := in.RefreshToken

	if token := refresh.ExtractTokenFromCookie(c); token != nil {
		refreshToken = token
	}

	if refreshToken == nil {
		return echo.NewHTTPError(http.StatusBadRequest, "no valid refreshToken provided")
	}

	if err := h.authService.Logout(c.Request().Context(), *refreshToken); err != nil {
		if errors.Is(err, auth.ErrSessionNotFound) {
			return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
		}

		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.NoContent(http.StatusOK)
}
