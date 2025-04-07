package app

import "github.com/moevm/nosql1h25-writer/backend/internal/api/common/mw"

func (app *App) AuthMW() *mw.AuthMW {
	if app.authMW != nil {
		return app.authMW
	}

	app.authMW = mw.NewAuthMW(app.AuthService())
	return app.authMW
}
