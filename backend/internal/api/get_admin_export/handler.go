package get_admin_export

import (
	"github.com/labstack/echo/v4"

	"github.com/moevm/nosql1h25-writer/backend/internal/api"
)

type handler struct{}

func New() api.Handler {
	return &handler{}
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
	return nil
}
