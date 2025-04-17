package post_auth_register

import (
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/moevm/nosql1h25-writer/backend/internal/api"
	"github.com/moevm/nosql1h25-writer/backend/internal/api/common/decorator"
	"github.com/moevm/nosql1h25-writer/backend/internal/api/common/refresh"
	"github.com/moevm/nosql1h25-writer/backend/internal/service/auth"
	"github.com/moevm/nosql1h25-writer/backend/internal/service/users"
)

type handler struct {
	authService auth.Service
}

func New(authService auth.Service) api.Handler {
	return decorator.NewBindAndValidate(&handler{
		authService: authService,
	})
}

type Request struct {
	DisplayName string `json:"displayName" validate:"required,min=3,max=64" example:"username"`
	Email       string `json:"email" validate:"required,email" format:"email" example:"test@gmail.com"`
	Password    string `json:"password" validate:"required,min=8,max=72" example:"Password123"`
}

type Response struct {
	ID           primitive.ObjectID `json:"id" validate:"required" example:"582ebf010936ac3ba5cd00e4"`
	AccessToken  string             `json:"accessToken" validate:"required" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"`
	RefreshToken uuid.UUID          `json:"refreshToken" validate:"required" example:"289abe45-5920-4366-a12a-875ddb422ace"`
}

// Handle - Register handler
//
//	@Summary		Register user
//	@Description	Create new user and return ID with `refresh` and `access` tokens
//	@Tags			auth
//	@Param			request	body	Request	true	"user credentials"
//	@Accept			json
//	@Produce		json
//	@Success		201	{object}	Response
//	@Failure		400	{object}	echo.HTTPError
//	@Failure		409	{object}	echo.HTTPError
//	@Failure		500	{object}	echo.HTTPError
//	@Router			/auth/register [post]
func (h *handler) Handle(c echo.Context, in Request) error {
	authData, err := h.authService.Register(c.Request().Context(), auth.RegisterIn{
		DisplayName: in.DisplayName,
		Email:       in.Email,
		Password:    in.Password,
	})
	if err != nil {
		if errors.Is(err, users.ErrUserAlreadyExists) {
			return echo.NewHTTPError(http.StatusConflict, err.Error())
		}

		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	c.SetCookie(&http.Cookie{
		Name:     refresh.RefreshToken,
		Value:    authData.Session.RefreshToken.String(),
		Expires:  authData.Session.ExpiresAt,
		Path:     api.AuthCookiePath,
		HttpOnly: true,
	})

	return c.JSON(http.StatusCreated, Response{
		ID:           authData.Session.UserID,
		AccessToken:  authData.AccessToken,
		RefreshToken: authData.Session.RefreshToken,
	})
}
