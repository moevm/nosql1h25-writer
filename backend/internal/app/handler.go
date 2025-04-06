package app

import (
	"github.com/moevm/nosql1h25-writer/backend/internal/api"
	"github.com/moevm/nosql1h25-writer/backend/internal/api/get_health"
	"github.com/moevm/nosql1h25-writer/backend/internal/api/post_auth_login"
	"github.com/moevm/nosql1h25-writer/backend/internal/api/post_auth_refresh"
)

func (app *App) GetHealthHandler() api.Handler {
	if app.getHealthHandler != nil {
		return app.getHealthHandler
	}

	app.getHealthHandler = get_health.New(app.OrdersCollection())
	return app.getHealthHandler
}

func (app *App) PostAuthLoginHandler() api.Handler {
	if app.postAuthLoginHandler != nil {
		return app.postAuthLoginHandler
	}

	app.postAuthLoginHandler = post_auth_login.New(app.AuthService())
	return app.postAuthLoginHandler
}

func (app *App) PostAuthRefreshHandler() api.Handler {
	if app.postAuthRefreshHandler != nil {
		return app.postAuthRefreshHandler
	}

	app.postAuthRefreshHandler = post_auth_refresh.New(app.AuthService())
	return app.postAuthRefreshHandler
}
