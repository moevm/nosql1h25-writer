package app

import "github.com/sv-tools/mongoifc"

func (app *App) MainDb() mongoifc.Database {
	if app.mainDb != nil {
		return app.mainDb
	}

	app.mainDb = app.mongoClient.Database("main")
	return app.mainDb
}

func (app *App) OrdersCollection() mongoifc.Collection {
	if app.ordersCollection != nil {
		return app.ordersCollection
	}

	app.ordersCollection = app.MainDb().Collection("orders")
	return app.ordersCollection
}
