package get_admin_stats

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/samber/lo"

	"github.com/moevm/nosql1h25-writer/backend/internal/api"
	"github.com/moevm/nosql1h25-writer/backend/internal/api/common/decorator"
	"github.com/moevm/nosql1h25-writer/backend/internal/service/stats"
)

type handler struct {
	statsService stats.Service
}

func New(statsService stats.Service) api.Handler {
	return decorator.NewBindAndValidate(&handler{statsService: statsService})
}

type Request struct {
	X   string `query:"x" validate:"oneof=newest oldest rich poor name_asc name_desc freelancer_rating_asc freelancer_rating_desc client_rating_asc client_rating_desc" example:"newest"`
	Y   string `query:"y" validate:"oneof=newest oldest rich poor name_asc name_desc freelancer_rating_asc freelancer_rating_desc client_rating_asc client_rating_desc" example:"newest"`
	Agg string `query:"agg" validate:"oneof=avg min max sum count" example:"count"`
}

type Point struct {
	X string  `json:"x"`
	Y float64 `json:"y"`
}

type Response struct {
	Points []Point `json:"points"`
}

// Handle - Return stats handler
//
//	@Summary		Return stats
//	@Description	Return stats
//	@Tags			admin
//	@Security		JWT
//	@Param			x	query	string		true	"X"	example(user_created_at)
//	@Param			y	query	string		true	"Y"	example(count)
//	@Param			agg	query	string	    true	"Aggregation Type" example(count)
//	@Produce		json
//	@Success		200	{object}	Response
//	@Failure		401	{object}	echo.HTTPError
//	@Failure		403	{object}	echo.HTTPError
//	@Failure		500	{object}	echo.HTTPError
//	@Router			/admin/stats [get]
func (h *handler) Handle(c echo.Context, in Request) error {
	out, err := h.statsService.Graph(c.Request().Context(), stats.GraphIn{
		X:       in.X,
		Y:       in.Y,
		AggType: stats.Aggregation(in.Agg),
	})
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, Response{Points: lo.Map(out, func(item stats.GraphOut, _ int) Point {
		return Point{
			X: item.X,
			Y: item.Y,
		}
	})})
}
