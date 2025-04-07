package app

import (
	"github.com/moevm/nosql1h25-writer/backend/internal/api"
	"github.com/moevm/nosql1h25-writer/backend/internal/api/get_health"
	"github.com/moevm/nosql1h25-writer/backend/internal/api/get_users"
	usersService "github.com/moevm/nosql1h25-writer/backend/internal/service/users"
)

func (app *App) GetHealthHandler() api.Handler {
	if app.getHealthHandler != nil {
		return app.getHealthHandler
	}

	app.getHealthHandler = get_health.New(app.OrdersCollection())
	return app.getHealthHandler
}

// GetUsersHandler возвращает синглтон обработчика для GET /users.
func (app *App) GetUsersHandler() api.Handler {
	if app.getUsersHandler != nil {
		return app.getUsersHandler
	}

	userService := usersService.New()

	app.getUsersHandler = get_users.New(userService)

	return app.getUsersHandler
}
