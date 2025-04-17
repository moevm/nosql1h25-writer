package post_auth_register

import (
	"github.com/labstack/echo/v4"

	"github.com/moevm/nosql1h25-writer/backend/internal/api"
	"github.com/moevm/nosql1h25-writer/backend/internal/api/common/decorator"
	"github.com/moevm/nosql1h25-writer/backend/internal/service/auth"
	"github.com/moevm/nosql1h25-writer/backend/internal/service/users"
)

type handler struct {
	usersService users.Service
	authService  auth.Service
}

func New(usersService users.Service, authService auth.Service) api.Handler {
	return decorator.NewBindAndValidate(&handler{
		usersService: usersService,
		authService:  authService,
	})
}

type Request struct {
}

type Response struct {
}

// Handle - Register handler
//
//	@Summary		Register user
//	@Description	Create new user and return ID with `refresh` and `access` tokens
//	@Tags			auth
//	@Param			request	body	Request	true	"user credentials"
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	Response
//	@Failure		400	{object}	echo.HTTPError
//	@Failure		409	{object}	echo.HTTPError
//	@Failure		500	{object}	echo.HTTPError
//	@Router			/auth/register [post]
func (h *handler) Handle(c echo.Context, in Request) error {
	return nil
}
