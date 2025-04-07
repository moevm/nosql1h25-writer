package app

import (
	"net/http"
	"os"

	validatorV10 "github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	log "github.com/sirupsen/logrus"
	echoSwagger "github.com/swaggo/echo-swagger"

	_ "github.com/moevm/nosql1h25-writer/backend/docs"
)

type CustomValidator struct {
	validator *validatorV10.Validate
}

func NewValidator() echo.Validator {
	return &CustomValidator{validator: validatorV10.New()}
}

func (cv *CustomValidator) Validate(i interface{}) error {
	if err := cv.validator.Struct(i); err != nil {

		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return nil
}

func (app *App) EchoHandler() *echo.Echo {
	if app.echoHandler != nil {
		return app.echoHandler
	}

	handler := echo.New()
	handler.Validator = NewValidator()
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

	handler.GET("/users", app.GetUsersHandler().Handle)
}

func setLogsFile() *os.File {
	file, err := os.OpenFile("/logs/requests.log", os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		log.Fatalf("v1 - setLogsFile - os.OpenFile: %v", err)
	}
	return file
}
