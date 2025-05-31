package get_admin_export

import (
	"fmt"
	"net/http"

	"github.com/jonboulle/clockwork"
	"github.com/labstack/echo/v4"

	"github.com/moevm/nosql1h25-writer/backend/internal/api"
	"github.com/moevm/nosql1h25-writer/backend/pkg/mongodb/mongotools"
)

type handler struct {
	mongoDumper mongotools.MongoDumper
	clock       clockwork.Clock
}

func New(mongoDumper mongotools.MongoDumper, clock clockwork.Clock) api.Handler {
	return &handler{
		mongoDumper: mongoDumper,
		clock:       clock,
	}
}

// Handle - Export mongodb state and return file handler
//
//	@Summary		Export mongodb state and return file
//	@Description	Export mongodb state and return file
//	@Tags			admin
//	@Security		JWT
//	@Produce		application/gzip
//	@Success		200
//	@Failure		401	{object}	echo.HTTPError
//	@Failure		403	{object}	echo.HTTPError
//	@Failure		500	{object}	echo.HTTPError
//	@Router			/admin/export [get]
func (h *handler) Handle(c echo.Context) error {
	filePath := fmt.Sprintf("tmp/dump_%d", h.clock.Now().Unix())

	if err := h.mongoDumper.Dump(filePath); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.Attachment(filePath, "dump.gzip")
}
