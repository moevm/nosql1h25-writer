package post_auth_refresh

import (
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

	"github.com/moevm/nosql1h25-writer/backend/internal/api"
	"github.com/moevm/nosql1h25-writer/backend/internal/api/common"
	"github.com/moevm/nosql1h25-writer/backend/internal/service/auth"
)

type handler struct {
	authService auth.Service
}

func New(authService auth.Service) api.Handler {
	return common.NewBindAndValidate(&handler{authService: authService})
}

type Request struct {
	RefreshToken *uuid.UUID `json:"refreshToken" example:"0e8f711e-b713-4869-b528-059a74311482"`
}

type Response struct {
	AccessToken  string    `json:"accessToken" validate:"required" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"`
	RefreshToken uuid.UUID `json:"refreshToken" validate:"required" example:"289abe45-5920-4366-a12a-875ddb422ace"`
}

// Handle - Refresh tokens handler
//
// @Summary		Refresh tokens
// @Description	Refresh `access` and `refresh` token pair. `refreshToken` can be passed in cookie
// @Tags			auth
// @Param			refreshToken	body	Request	false	"active refresh token in UUID RFC4122 format"
// @Accept			json
// @Produce		json
// @Success		200	{object}	Response
// @Failure		400	{object}	echo.HTTPError
// @Failure		401	{object}	echo.HTTPError
// @Failure		500	{object}	echo.HTTPError
// @Router			/auth/refresh [post]
func (h *handler) Handle(c echo.Context, in Request) error {
	refreshToken := in.RefreshToken

	if token := common.ExtractRefreshTokenFromCookie(c); token != nil {
		refreshToken = token
	}

	if refreshToken == nil {
		return echo.NewHTTPError(http.StatusBadRequest, "no valid refreshToken provided")
	}

	authData, err := h.authService.Refresh(c.Request().Context(), *refreshToken)
	if err != nil {
		if errors.Is(err, auth.ErrSessionNotFound) || errors.Is(err, auth.ErrSessionExpired) {
			return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
		}

		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	c.SetCookie(&http.Cookie{
		Name:     api.RefreshToken,
		Value:    authData.Session.RefreshToken.String(),
		Expires:  authData.Session.ExpiresAt,
		Path:     api.AuthCookiePath,
		HttpOnly: true,
	})
	c.Path()
	return c.JSON(http.StatusOK, Response{
		AccessToken:  authData.AccessToken,
		RefreshToken: authData.Session.RefreshToken,
	})
}
