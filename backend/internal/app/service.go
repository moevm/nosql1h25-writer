package app

import (
	"github.com/moevm/nosql1h25-writer/backend/internal/service/auth"
	"github.com/moevm/nosql1h25-writer/backend/internal/service/orders"
	"github.com/moevm/nosql1h25-writer/backend/internal/service/users"
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

func (app *App) OrdersService() orders.Service {
	if app.ordersService != nil {
		return app.ordersService
	}

	app.ordersService = orders.New(
		app.OrdersRepo(),
		app.UsersService(),
	)
	return app.ordersService
}