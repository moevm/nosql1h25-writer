package get_health

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/moevm/nosql1h25-writer/backend/api"
	"github.com/sv-tools/mongoifc"
)

type handler struct {
	orders mongoifc.Collection
}

func New(orders mongoifc.Collection) api.Handler {
	return &handler{orders: orders}
}

func (h *handler) Handle(c echo.Context) error {
	return c.String(http.StatusOK, h.orders.Name())
}
