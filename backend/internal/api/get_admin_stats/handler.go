package get_admin_stats

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/samber/lo"

	"github.com/moevm/nosql1h25-writer/backend/internal/api"
	"github.com/moevm/nosql1h25-writer/backend/internal/api/common/decorator"
	"github.com/moevm/nosql1h25-writer/backend/internal/entity"
	"github.com/moevm/nosql1h25-writer/backend/internal/service/stats"
)

type handler struct {
	statsService stats.Service
}

func New(statsService stats.Service) api.Handler {
	return decorator.NewBindAndValidate(&handler{statsService: statsService})
}

type Request struct {
	X   string `query:"x" validate:"oneof=user_id user_system_role user_active user_created_at order_id order_active order_freelancer_id order_client_id order_created_at" example:"newest"`
	Y   string `query:"y" validate:"oneof=count user_balance user_client_rating user_freelancer_rating order_completion_time order_cost order_responses_count" example:"newest"`
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
//	@Param			x	query	string	true	"X"					example(user_created_at)
//	@Param			y	query	string	true	"Y"					example(count)
//	@Param			agg	query	string	true	"Aggregation Type"	example(count)
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
		AggType: entity.Aggregation(in.Agg),
	})
	if err != nil {
		if errors.Is(err, stats.ErrInvalidRequest) {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, Response{Points: lo.Map(out, func(item stats.GraphOut, _ int) Point {
		return Point{
			X: item.X,
			Y: item.Y,
		}
	})})
}
