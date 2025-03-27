# Написание обработчика
## Инструкция
1. Создайте в [`/internal/api`](../internal/api) каталог и назовите его в формате "метод_путь_до_страницы". Например: `GET /admin/export` → `get_admin_export`, `PATCH /orders/:id` → `patch_orders_id`.
2. Создайте в этом каталоге файл `handler.go`.
3. Не забудьте заполнить документацию для **Swagger** (комментарии перед функцией `Handle`).

* Документация по заполнению Swagger: [ссылка](https://github.com/swaggo/swag/blob/master/README.md).

## Валидация входных данных
* Объявите структуру в файле `handler.go`, например (*string используется, чтобы указать, что поле является необязательным):
```go
type updateSongInput struct {
	Id          int     `param:"id" validate:"number,gt=0"`
	Title       *string `json:"group" validate:"omitempty,max=128" example:"Promo for vegetables"`
	Song        *string `json:"song" validate:"omitempty,max=128" example:"Best Compilation"`
	Link        *string `json:"link" validate:"omitempty,max=128,uri" example:"https://www.youtube.com/watch?v=Xsp3_a-PMTw"`
	ReleaseDate *string `json:"releaseDate" validate:"omitempty,date" example:"2006-06-22"`
}
```

* Документация по валидатору: [ссылка](https://github.com/go-playground/validator/blob/master/README.md).

* Далее внутри метода `Handle` используйте следующий код:
```go
	var input updateSongInput

	if err := c.Bind(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid request")
		return err
	}

	if err := c.Validate(input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return err
	}
```

## Образец `handler.go`
```go
package get_users

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/sv-tools/mongoifc"

	"github.com/moevm/nosql1h25-writer/backend/internal/api"
)

type handler struct {
	users users.Service
}

func New(users users.Service) api.Handler {
	return &handler{users: users}
}

// @Description Получить пользователя с ID 1
// @Summary Получить пользователя
// @Param id path int true "ID пользователя" minimum(1) example(33)
// @Produce json
// @Success 200 {object} entity.User
// @Failure 400 {object} echo.HTTPError
// @Failure 500 {object} echo.HTTPError
// @Router /users [get]
func (h *handler) Handle(c echo.Context) error {
	return c.String(http.StatusOK, h.users.GetUser(1))
}
```

## Пример документации
```go
// @Description Поиск пользователей с сортировкой и пагинацией
// @Summary Поиск пользователей
// @Param order_by query string false "Список критериев для сортировки. Если не указан порядок, то по возрастанию." example(id:asc,name:desc,created_at)
// @Param offset query int false "Offset" default(0) minimum(0) example(10)
// @Param limit query int false "Limit" default(10) minimum(1) maximum(20) example(10)
// @Produce json
// @Success 200 {array} entity.User
// @Failure 400 {object} echo.HTTPError
// @Failure 500 {object} echo.HTTPError
// @Router /users [get]
```

```go
// @Description Удалить пользователя по ID
// @Summary Удалить пользователя
// @Param id path int true "ID пользователя" minimum(1) example(33)
// @Success 204
// @Failure 400 {object} echo.HTTPError
// @Failure 500 {object} echo.HTTPError
// @Router /users/{id} [delete]
```

## Название маршрутов
Гайд как давать названия маршрутам: [ссылка](https://restfulapi.net/resource-naming/).