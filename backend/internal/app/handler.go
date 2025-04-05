package app

import (
	"github.com/moevm/nosql1h25-writer/backend/internal/api"
	"github.com/moevm/nosql1h25-writer/backend/internal/api/get_health"
	"github.com/moevm/nosql1h25-writer/backend/internal/api/get_orders"
	"github.com/moevm/nosql1h25-writer/backend/internal/api/get_orders_id"
)

func (app *App) GetHealthHandler() api.Handler {
	if app.getHealthHandler != nil {
		return app.getHealthHandler
	}

	app.getHealthHandler = get_health.New(app.OrdersCollection())
	return app.getHealthHandler
}

func (app *App) GetOrdersHandler() api.Handler {
	if app.getOrdersHandler != nil {
		return app.getOrdersHandler
	}

	app.getOrdersHandler = get_orders.New(app.OrdersCollection())
	return app.getOrdersHandler
}

func (app *App) GetOrdersIDHandler() api.Handler {
	if app.getOrdersIDHandler != nil {
		return app.getOrdersIDHandler
	}

	app.getOrdersIDHandler = get_orders_id.New(app.OrdersCollection())
	return app.getOrdersIDHandler
}
