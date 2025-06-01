package app

import (
	"github.com/jonboulle/clockwork"

	"github.com/moevm/nosql1h25-writer/backend/pkg/hasher"
	"github.com/moevm/nosql1h25-writer/backend/pkg/mongodb/mongotools"
)

func (app *App) PasswordHasher() hasher.PasswordHasher {
	if app.passwordHasher != nil {
		return app.passwordHasher
	}

	app.passwordHasher = hasher.NewBcrypt()
	return app.passwordHasher
}

func (app *App) Clock() clockwork.Clock {
	if app.clock != nil {
		return app.clock
	}

	app.clock = clockwork.NewRealClock()
	return app.clock
}

func (app *App) MongoDumper() mongotools.MongoDumper {
	if app.mongoDumper != nil {
		return app.mongoDumper
	}

	app.mongoDumper = mongotools.NewDumper(app.cfg.Mongo.Uri)
	return app.mongoDumper
}
