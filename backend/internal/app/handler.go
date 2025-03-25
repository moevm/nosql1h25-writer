package app

import (
	"github.com/moevm/nosql1h25-writer/backend/internal/api"
	"github.com/moevm/nosql1h25-writer/backend/internal/api/get_health"
)

func (app *App) getGetHealthHandler() api.Handler {
	if app.getHealthHandler != nil {
		return app.getHealthHandler
	}

	app.getHealthHandler = get_health.New(app.getOrdersCollection())
	return app.getHealthHandler
}
