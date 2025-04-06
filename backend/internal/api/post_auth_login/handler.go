package post_auth_login

import (
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

	"github.com/moevm/nosql1h25-writer/backend/internal/api"
	"github.com/moevm/nosql1h25-writer/backend/internal/api/common"
	"github.com/moevm/nosql1h25-writer/backend/internal/service/auth"
	"github.com/moevm/nosql1h25-writer/backend/internal/service/users"
)

type handler struct {
	authService auth.Service
}

func New(authService auth.Service) api.Handler {
	return common.NewBindAndValidate(&handler{authService: authService})
}

type Request struct {
	Email    string `json:"email" validate:"required,email" format:"email" example:"test@gmail.com"`
	Password string `json:"password" validate:"required,min=8,max=72" minLength:"8" maxLength:"72" example:"Password123"`
}

type Response struct {
	AccessToken  string    `json:"accessToken" validate:"required" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"`
	RefreshToken uuid.UUID `json:"refreshToken" validate:"required" example:"289abe45-5920-4366-a12a-875ddb422ace"`
}

// @Summary		Login by email and password
// @Description	Generate `access` and `refresh` token pair. `refreshToken` sets in httpOnly cookie also.
// @Tags auth
// @Param		request	body	Request	true	"existing user credentials"
// @Accept			json
// @Produce		json
// @Success		200	{object}	Response
// @Failure		400	{object}	echo.HTTPError
// @Failure	404	{object}	echo.HTTPError
// @Failure		500	{object}	echo.HTTPError
// @Router			/auth/login [post]
func (h *handler) Handle(c echo.Context, in Request) error {
	authData, err := h.authService.Login(c.Request().Context(), in.Email, in.Password)
	if err != nil {
		if errors.Is(err, users.ErrUserNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, err.Error())
		} else if errors.Is(err, auth.ErrWrongPassword) {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	c.SetCookie(&http.Cookie{
		Name:     "refreshToken",
		Value:    authData.Session.RefreshToken.String(),
		Expires:  authData.Session.ExpiresAt,
		Path:     "/api/auth",
		HttpOnly: true,
	})

	return c.JSON(http.StatusOK, Response{
		AccessToken:  authData.AccessToken,
		RefreshToken: authData.Session.RefreshToken,
	})
}
