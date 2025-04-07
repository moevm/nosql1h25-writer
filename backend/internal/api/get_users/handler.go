package get_users

import (
	"context"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"

	"github.com/moevm/nosql1h25-writer/backend/internal/api"
	"github.com/moevm/nosql1h25-writer/backend/internal/entity"
)

type userService interface {
	FindUsers(ctx context.Context, params entity.UserSearchParams) ([]entity.User, int64, error)
}

type handler struct {
	users userService
}

func New(users userService) api.Handler {
	return &handler{users: users}
}

type input struct {
	Offset  int64    `query:"offset" validate:"omitempty,number,min=0" default:"0"`
	Limit   int64    `query:"limit" validate:"omitempty,number,min=1,max=50" default:"10"`
	Profile []string `query:"profile" validate:"omitempty,dive,oneof=client freelancer"`
}

type profileOutput struct {
	Role        string  `json:"role"`
	Description *string `json:"description,omitempty"`
	Rating      float64 `json:"rating"`
}

type userOutput struct {
	ID          string          `json:"id"`
	DisplayName string          `json:"displayName"`
	Profiles    []profileOutput `json:"profiles"`
	CreatedAt   string          `json:"createdAt"`
	UpdatedAt   string          `json:"updatedAt"`
}

type output struct {
	Total int64        `json:"total"`
	Users []userOutput `json:"users"`
}

// Handle - основной метод обработчика, который вызывается при получении запроса GET /users.
// @Summary      Получить список пользователей
// @Description  Возвращает список пользователей с пагинацией и возможностью фильтрации по типу профиля.
// @Tags         Users
// @Accept       json
// @Produce      json
// @Param        offset    query     int      false  "Смещение для пагинации (сколько записей пропустить)" default(0) minimum(0)
// @Param        limit     query     int      false  "Количество записей на страницу" default(10) minimum(1) maximum(50)
// @Param        profile   query     []string false  "Фильтр по типу профиля ('client', 'freelancer'). Можно указать несколько раз." Enums(client, freelancer) collectionFormat(multi)
// @Success      200       {object}  output   "Успешный ответ со списком пользователей и общим количеством"
// @Failure      400       {object}  echo.HTTPError "Ошибка валидации входных данных"
// @Failure      500       {object}  echo.HTTPError "Внутренняя ошибка сервера"
// @Router       /users [get]

func (h *handler) Handle(c echo.Context) error {
	var inp input

	if err := c.Bind(&inp); err != nil {

		return echo.NewHTTPError(http.StatusBadRequest, "Invalid input parameters: "+err.Error())
	}

	if err := c.Validate(inp); err != nil {

		return err
	}

	if c.QueryParam("offset") == "" {
		inp.Offset = 0
	}
	if c.QueryParam("limit") == "" {
		inp.Limit = 10
	}

	searchParams := entity.UserSearchParams{
		Offset:        inp.Offset,
		Limit:         inp.Limit,
		ProfileFilter: inp.Profile,
	}

	ctx := c.Request().Context()
	users, total, err := h.users.FindUsers(ctx, searchParams)
	if err != nil {

		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to fetch users")
	}

	usersOutput := make([]userOutput, 0, len(users))
	for _, u := range users {
		profilesOut := make([]profileOutput, 0, len(u.Profiles))
		for _, p := range u.Profiles {
			profilesOut = append(profilesOut, profileOutput{
				Role:        p.Role,
				Description: p.Description,
				Rating:      p.Rating,
			})
		}

		usersOutput = append(usersOutput, userOutput{
			ID:          u.ID.Hex(),
			DisplayName: u.DisplayName,
			Profiles:    profilesOut,
			CreatedAt:   u.CreatedAt.Format(time.RFC3339),
			UpdatedAt:   u.UpdatedAt.Format(time.RFC3339),
		})
	}

	response := output{
		Total: total,
		Users: usersOutput,
	}
	return c.JSON(http.StatusOK, response)
}
