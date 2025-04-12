package app

import (
	"github.com/jonboulle/clockwork"

	"github.com/moevm/nosql1h25-writer/backend/pkg/hasher"
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
