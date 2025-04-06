package api

import "github.com/labstack/echo/v4"

const (
	AuthCookiePath = "/api/auth"
	RefreshToken   = "refreshToken"
)

type Handler interface {
	Handle(c echo.Context) error
}
