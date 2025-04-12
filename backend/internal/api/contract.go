package api

import "github.com/labstack/echo/v4"

const AuthCookiePath = "/api/auth"

type Handler interface {
	Handle(c echo.Context) error
}
