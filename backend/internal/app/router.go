package app

import (
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	log "github.com/sirupsen/logrus"
	echoSwagger "github.com/swaggo/echo-swagger"

	_ "github.com/moevm/nosql1h25-writer/backend/docs"
	"github.com/moevm/nosql1h25-writer/backend/pkg/validator"
)

func (app *App) EchoHandler() *echo.Echo {
	if app.echoHandler != nil {
		return app.echoHandler
	}

	handler := echo.New()
	handler.Validator = validator.NewCustomValidator()
	handler.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: `{"time":"${time_rfc3339_nano}", "method":"${method}","uri":"${uri}", "status":${status},"error":"${error}"}` + "\n",
		Output: setLogsFile(),
	}))
	handler.Use(middleware.Recover())
	handler.GET("/swagger/*", echoSwagger.WrapHandler)

	app.configureRouter(handler)

	app.echoHandler = handler
	return app.echoHandler
}

func (app *App) configureRouter(handler *echo.Echo) {
	handler.GET("/health", app.GetHealthHandler().Handle)

	authGroup := handler.Group("/auth")
	{
		authGroup.POST("/login", app.PostAuthLoginHandler().Handle)
		authGroup.POST("/refresh", app.PostAuthRefreshHandler().Handle)
		authGroup.POST("/logout", app.PostAuthLogoutHandler().Handle)
	}

	adminGroup := handler.Group("/admin", app.AuthMW().UserIdentity())
	{
		adminGroup.GET("", app.GetAdminHandler().Handle, app.AuthMW().AdminRole())
	}

	balanceGroup := handler.Group("/balance", app.AuthMW().UserIdentity())
	{
		balanceGroup.POST("/deposit", app.PostBalanceDepositHandler().Handle)
		balanceGroup.POST("/withdraw", app.PostBalanceWithdrawHandler().Handle)
	}
}

func setLogsFile() *os.File {
	file, err := os.OpenFile("/logs/requests.log", os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		log.Fatalf("v1 - setLogsFile - os.OpenFile: %v", err)
	}
	return file
}
