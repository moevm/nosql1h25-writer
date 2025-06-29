package app

import (
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	log "github.com/sirupsen/logrus"
	echoSwagger "github.com/swaggo/echo-swagger"

	_ "github.com/moevm/nosql1h25-writer/backend/docs"
	"github.com/moevm/nosql1h25-writer/backend/internal/entity"
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

	usersGroup := handler.Group("/users")
	{
		usersGroup.GET("/:id", app.GetUsersIDHandler().Handle, app.AuthMW().UserIdentity())
		usersGroup.PATCH("/:id", app.PatchUsersIDHandler().Handle, app.AuthMW().UserIdentity())
		usersGroup.GET("/:id/orders", app.GetUsersIDOrdersHandler().Handle, app.AuthMW().UserIdentity())
		usersGroup.GET("/:id/responses", app.GetUsersIDResponsesHandler().Handle, app.AuthMW().UserIdentity())
	}

	authGroup := handler.Group("/auth")
	{
		authGroup.POST("/register", app.PostAuthRegisterHandler().Handle)
		authGroup.POST("/login", app.PostAuthLoginHandler().Handle)
		authGroup.POST("/refresh", app.PostAuthRefreshHandler().Handle)
		authGroup.POST("/logout", app.PostAuthLogoutHandler().Handle)
	}

	adminGroup := handler.Group("/admin", app.AuthMW().UserIdentity())
	{
		adminGroup.GET("", app.GetAdminHandler().Handle, app.AuthMW().Role(entity.SystemRoleTypeAdmin))
		adminGroup.GET("/export", app.GetAdminExportHandler().Handle, app.AuthMW().Role(entity.SystemRoleTypeAdmin))
		adminGroup.POST("/import", app.PostAdminImportHandler().Handle, app.AuthMW().Role(entity.SystemRoleTypeAdmin))
		adminGroup.GET("/users", app.GetAdminUsersHandler().Handle, app.AuthMW().Role(entity.SystemRoleTypeAdmin))
		adminGroup.GET("/stats", app.GetAdminStatsHandler().Handle, app.AuthMW().Role(entity.SystemRoleTypeAdmin))
	}

	balanceGroup := handler.Group("/balance", app.AuthMW().UserIdentity())
	{
		balanceGroup.POST("/deposit", app.PostBalanceDepositHandler().Handle)
		balanceGroup.POST("/withdraw", app.PostBalanceWithdrawHandler().Handle)
	}

	ordersGroup := handler.Group("/orders", app.authMW.UserIdentity())
	{
		ordersGroup.POST("", app.PostOrdersHandler().Handle)
		ordersGroup.GET("", app.GetOrdersHandler().Handle)
		ordersGroup.GET("/:id", app.GetOrdersIDHandler().Handle)
		ordersGroup.PATCH("/:id", app.PatchOrdersIDHandler().Handle)
		ordersGroup.POST("/:id/response", app.PostOrdersResponseHandler().Handle)
	}
}

func setLogsFile() *os.File {
	file, err := os.OpenFile("/logs/requests.log", os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		log.Fatalf("v1 - setLogsFile - os.OpenFile: %v", err)
	}
	return file
}
