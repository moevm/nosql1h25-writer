package app

import "github.com/sv-tools/mongoifc"

func (app *App) getMainDb() mongoifc.Database {
	if app.mainDb != nil {
		return app.mainDb
	}

	app.mainDb = app.mongoClient.Database("main")
	return app.mainDb
}

func (app *App) getOrdersCollection() mongoifc.Collection {
	if app.ordersCollection != nil {
		return app.ordersCollection
	}

	app.ordersCollection = app.getMainDb().Collection("orders")
	return app.ordersCollection
}
