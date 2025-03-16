package app

import (
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/moevm/nosql1h25-writer/backend/api/get_health"
	log "github.com/sirupsen/logrus"
)

func (app *App) getEchoHandler() *echo.Echo {
	if app.echoHandler != nil {
		return app.echoHandler
	}

	handler := echo.New()
	handler.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: `{"time":"${time_rfc3339_nano}", "method":"${method}","uri":"${uri}", "status":${status},"error":"${error}"}` + "\n",
		Output: setLogsFile(),
	}))
	handler.Use(middleware.Recover())
	app.configureRouter(handler)

	app.echoHandler = handler
	return app.echoHandler
}

func (app *App) configureRouter(handler *echo.Echo) {
	apiGroup := handler.Group("/api")
	{
		apiGroup.GET("/health", get_health.New(app.getOrdersCollection()).Handle)
	}
}

func setLogsFile() *os.File {
	file, err := os.OpenFile("/logs/requests.log", os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		log.Fatalf("v1 - setLogsFile - os.OpenFile: %v", err)
	}
	return file
}
