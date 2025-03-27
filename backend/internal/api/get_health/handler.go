package get_health

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/sv-tools/mongoifc"

	"github.com/moevm/nosql1h25-writer/backend/internal/api"
)

type handler struct {
	orders mongoifc.Collection
}

func New(orders mongoifc.Collection) api.Handler {
	return &handler{orders: orders}
}

//	@Summary		Check health
//	@Description	Whether REST-API alive or not
//	@Tags			health
//	@Produce		plain
//	@Success		200	{string}	string
//	@Failure		500	{object}	echo.HTTPError
//	@Router			/health [get]
func (h *handler) Handle(c echo.Context) error {
	return c.String(http.StatusOK, h.orders.Name())
}
