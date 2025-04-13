package get_users

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"

	"github.com/moevm/nosql1h25-writer/backend/internal/api"
	"github.com/moevm/nosql1h25-writer/backend/internal/api/common/decorator"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/moevm/nosql1h25-writer/backend/internal/entity"
	usersExt "github.com/moevm/nosql1h25-writer/backend/internal/service/usersExt"
)

const defaultLimit = 10

type handler struct {
	users usersExt.Service
}

func New(users usersExt.Service) api.Handler {
	return decorator.NewBindAndValidate(&handler{users: users})
}

type Request struct {
	Offset  int64    `query:"offset" validate:"omitempty,number,min=0" default:"0"`
	Limit   int64    `query:"limit" validate:"omitempty,number,min=1,max=50" default:"10"`
	Profile []string `query:"profile" validate:"omitempty,dive,oneof=client freelancer"`
}

type profileOut struct {
	Role        string  `json:"role"`
	Description *string `json:"description,omitempty"`
	Rating      float64 `json:"rating"`
}

type userOut struct {
	ID          primitive.ObjectID `json:"id"`
	DisplayName string             `json:"displayName"`
	Profiles    []profileOut       `json:"profiles"`
	CreatedAt   time.Time          `json:"createdAt"`
	UpdatedAt   time.Time          `json:"updatedAt"`
}

type Response struct {
	Total int64     `json:"total"`
	Users []userOut `json:"users"`
}

// Handle - основной метод обработчика, который вызывается при получении запроса GET /users.
//
//	@Summary		Получить список пользователей
//	@Description	Возвращает список пользователей с пагинацией и возможностью фильтрации по типу профиля.
//	@Tags			Users
//	@Accept			json
//	@Produce		json
//	@Param			offset	query		int64			false	"Смещение для пагинации (сколько записей пропустить)"							default(0)					minimum(0)
//	@Param			limit	query		int64			false	"Количество записей на страницу"												default(10)					minimum(1)	maximum(50)
//	@Param			profile	query		[]string		false	"Фильтр по типу профиля ('client', 'freelancer'). Можно указать несколько раз."	Enums(client, freelancer)	collectionFormat(multi)
//	@Success		200		{object}	Response		"Успешный ответ со списком пользователей и общим количеством"
//	@Failure		400		{object}	echo.HTTPError	"Ошибка валидации входных данных"
//	@Failure		500		{object}	echo.HTTPError	"Внутренняя ошибка сервера"
//	@Router			/users [get]
func (h *handler) Handle(c echo.Context, inp Request) error {
	if c.QueryParam("offset") == "" {
		inp.Offset = 0
	}
	if c.QueryParam("limit") == "" {
		inp.Limit = defaultLimit
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

	usersOutput := make([]userOut, 0, len(users))
	for _, u := range users {
		userExt := u
		profilesOut := make([]profileOut, 0, len(userExt.Profiles))
		for _, p := range userExt.Profiles {
			profilesOut = append(profilesOut, profileOut{
				Role:        p.Role,
				Description: &p.Description,
				Rating:      p.Rating,
			})
		}

		usersOutput = append(usersOutput, userOut{
			ID:          userExt.ID,
			DisplayName: userExt.DisplayName,
			Profiles:    profilesOut,
			CreatedAt:   userExt.CreatedAt,
			UpdatedAt:   userExt.UpdatedAt,
		})
	}

	response := Response{
		Total: total,
		Users: usersOutput,
	}
	return c.JSON(http.StatusOK, response)
}