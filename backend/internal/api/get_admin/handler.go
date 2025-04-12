package get_admin

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/moevm/nosql1h25-writer/backend/internal/api"
	"github.com/moevm/nosql1h25-writer/backend/internal/api/common/mw"
	"github.com/moevm/nosql1h25-writer/backend/internal/entity"
)

type handler struct{}

func New() api.Handler {
	return &handler{}
}

type Response struct {
	SystemRole entity.SystemRoleType `json:"systemRole" validate:"required" example:"admin"`
	UserID     primitive.ObjectID    `json:"userId" validate:"required" example:"5a2493c33c95a1281836eb6a"`
}

// Handle - Check admin rights available handler
//
//	@Summary		Check admin rights available
//	@Description	Whether user has admin rights
//	@Tags			admin
//	@Security		JWT
//	@Produce		json
//	@Success		200	{object}	Response
//	@Failure		500	{object}	echo.HTTPError
//	@Router			/admin [get]
func (h *handler) Handle(c echo.Context) error {
	return c.JSON(http.StatusOK, Response{
		SystemRole: c.Get(mw.SystemRoleKey).(entity.SystemRoleType), //nolint:forcetypeassert
		UserID:     c.Get(mw.UserIDKey).(primitive.ObjectID),        //nolint:forcetypeassert
	})
}
