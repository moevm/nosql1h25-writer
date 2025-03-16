package get_health

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/moevm/nosql1h25-writer/backend/api"
)

type handler struct{}

func New() api.Handler {
	return &handler{}
}

func (h *handler) Handle(c echo.Context) error {
	return c.String(http.StatusOK, "OK")
}
