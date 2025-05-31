package post_admin_import

import (
	"fmt"
	"io"
	"net/http"
	"os"

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

// Handle - Import mongodb state into current instance handler
//
//	@Summary		Import mongodb state into current instance
//	@Description	Import mongodb state into current instance
//	@Tags			admin
//	@Security		JWT
//	@Accept			application/gzip
//	@Produce		json
//	@Success		200
//	@Failure		401	{object}	echo.HTTPError
//	@Failure		403	{object}	echo.HTTPError
//	@Failure		500	{object}	echo.HTTPError
//	@Router			/admin/import [post]
func (h *handler) Handle(c echo.Context) error {
	inFormFile, err := c.FormFile("file")
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "no file provided")
	}

	src, err := inFormFile.Open()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "cannot read file")
	}
	defer src.Close() //nolint:errcheck

	filePath := fmt.Sprintf("tmp/restore_%d", h.clock.Now().Unix())
	dst, err := os.Create(filePath)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "cannot create file")
	}
	defer dst.Close() //nolint:errcheck

	if _, err := io.Copy(dst, src); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "cannot copy content")
	}

	if err := h.mongoDumper.Restore(filePath); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "successful import"})
}
