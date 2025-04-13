package app

import (
	"github.com/moevm/nosql1h25-writer/backend/internal/service/auth"
	users "github.com/moevm/nosql1h25-writer/backend/internal/service/users"
	users_ext_service "github.com/moevm/nosql1h25-writer/backend/internal/service/usersExt"
)

func (app *App) AuthService() auth.Service {
	if app.authService != nil {
		return app.authService
	}

	app.authService = auth.New(
		app.UsersRepo(),
		app.AuthRepo(),
		app.PasswordHasher(),
		app.Clock(),
		app.cfg.Auth.JWTSecretKey,
		app.cfg.Auth.AccessTokenTTL,
		app.cfg.Auth.RefreshTokenTTL,
	)
	return app.authService
}

// UsersService возвращает синглтон сервиса users.
func (app *App) UsersService() users.Service {
	if app.usersService != nil {
		return app.usersService
	}

	app.usersService = users.New(
		app.UsersRepo(),
	)
	return app.usersService
}

// UsersExtService возвращает синглтон сервиса usersExt.
func (app *App) UsersExtService() users_ext_service.Service {
	if app.usersExtService != nil {
		return app.usersExtService
	}
	app.usersExtService = users_ext_service.New()

	return app.usersExtService
}