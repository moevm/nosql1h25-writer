package app

import "github.com/moevm/nosql1h25-writer/backend/internal/service/auth"

func (app *App) AuthService() auth.Service {
	if app.authService != nil {
		return app.authService
	}

	app.authService = auth.New(
		app.UsersRepo(),
		app.AuthRepo(),
		app.PasswordHasher(),
		app.Clock(),
		app.cfg.Auth.JWTSecretKey,
		app.cfg.Auth.AccessTokenTTL,
		app.cfg.Auth.RefreshTokenTTL,
	)
	return app.authService
}
