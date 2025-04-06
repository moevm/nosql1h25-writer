package app

import (
	"github.com/sv-tools/mongoifc"

	"github.com/moevm/nosql1h25-writer/backend/internal/repo/auth"
	"github.com/moevm/nosql1h25-writer/backend/internal/repo/users"
)

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

func (app *App) UsersCollection() mongoifc.Collection {
	if app.usersCollection != nil {
		return app.usersCollection
	}

	app.usersCollection = app.MainDb().Collection("users")
	return app.usersCollection
}

func (app *App) SessionsCollection() mongoifc.Collection {
	if app.sessionsCollection != nil {
		return app.sessionsCollection
	}

	app.sessionsCollection = app.MainDb().Collection("sessions")
	return app.sessionsCollection
}

func (app *App) UsersRepo() users.Repo {
	if app.usersRepo != nil {
		return app.usersRepo
	}

	app.usersRepo = users.New(app.UsersCollection())
	return app.usersRepo
}

func (app *App) AuthRepo() auth.Repo {
	if app.authRepo != nil {
		return app.authRepo
	}

	app.authRepo = auth.New(app.SessionsCollection())
	return app.authRepo
}
