package app

import (
	users_ext_repo "github.com/moevm/nosql1h25-writer/backend/internal/repo/usersExt"
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

func (app *App) UsersService() users.Service {
	if app.usersService != nil {
		return app.usersService
	}

	app.usersService = users.New(
		app.UsersRepo(),
	)
	return app.usersService
}

func (app *App) UsersExtService() users_ext_service.Service {
	if app.usersExtService != nil {
		return app.usersExtService
	}

	repo := users_ext_repo.New(app.MainDb())

	app.usersExtService = users_ext_service.New(repo)

	return app.usersExtService
}
