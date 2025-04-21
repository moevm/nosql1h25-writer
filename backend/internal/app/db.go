package app

import (
	"context"

	"github.com/sv-tools/mongoifc"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/moevm/nosql1h25-writer/backend/internal/repo/auth"
	"github.com/moevm/nosql1h25-writer/backend/internal/repo/orders"
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
	_, err := app.usersCollection.Indexes().CreateOne(context.Background(), mongo.IndexModel{
		Keys:    bson.M{"email": -1},
		Options: options.Index().SetUnique(true),
	})
	if err != nil {
		panic("invalid index setup")
	}

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

	app.usersRepo = users.New(app.UsersCollection(), app.Clock())
	return app.usersRepo
}

func (app *App) AuthRepo() auth.Repo {
	if app.authRepo != nil {
		return app.authRepo
	}

	app.authRepo = auth.New(app.SessionsCollection(), app.Clock())
	return app.authRepo
}

func (app *App) OrdersRepo() orders.Repo {
	if app.ordersRepo != nil {
		return app.ordersRepo
	}

	app.ordersRepo = orders.New(app.OrdersCollection(), app.Clock())
	return app.ordersRepo
}
