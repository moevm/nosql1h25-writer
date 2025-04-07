# Написание обработчика
## Обзор архитектуры проекта
Проект организован по принципам чистой архитектуры, где разделены слои:
* API (ручки): принимают входящие HTTP-запросы, валидируют данные, вызывает соответствующие сервисы и формируют ответ. Их код находится в папке [`internal/api`](../internal/api).
* Слой сервисов: содержит бизнес-логику (аутентификация, работа с заказами, пользователями). Располагается в папке [`internal/service`](../internal/service).
* Репозиторий: отвечает за взаимодействие с базой данных или другими источниками данных. Располагается в папке [`internal/repo`](../internal/repo).
* Общие компоненты: папка [`internal/api/common`](../internal/api/common) содержит общие элементы (декораторы, middleware, работа с refresh-токенами).

## Создание нового обработчика
При разработке нового эндпоинта следуйте по этому алгоритму:
1. Создайте новую папку в [`internal/api`](../internal/api) с именем, соответствующим функциональности (название метода + маршрут, например: `PATCH /orders/:id` → `patch_orders_id`).
   * Как давать названия маршрутам, лучшие практики: [ссылка](https://restfulapi.net/resource-naming/).
2. Создайте в этой папке файл `handler.go` и вставьте следующий код (обратите внимание на **комментарии**):
```go
package post_example

import (
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

	"github.com/moevm/nosql1h25-writer/backend/internal/api"
	"github.com/moevm/nosql1h25-writer/backend/internal/api/common/decorator"
	"github.com/moevm/nosql1h25-writer/backend/internal/service/example"
)

type handler struct {
    // Здесь перечислите все сервисы, которые вам нужны (обычно нужен лишь один сервис и один вызов метода).
    // Например: authService auth.Service.
    exampleService example.Service
}

// Не забудьте указать нужный сервис здесь (а этот комментарий стереть).
func New(exampleService example.Service) api.Handler {
	return decorator.NewBindAndValidate(&handler{exampleService: exampleService})
}

// Что должно быть в запросе.
// Для каждого поля нужно заполнить примеры для документации и инструкции для валидатора.
type Request struct {
    // Пример заполнения.
    // Ссылки на дополнительные ресурсы:
    // validate: https://github.com/go-playground/validator/blob/master/README.md
    // format, example, ...: https://github.com/swaggo/swag/blob/master/README.md 
	ID       string `param:"id" validate:"required,alnum" example:"m19lfjDkdffm"` // Это path-параметр (/users/:id/).
	Email    string `json:"email" validate:"required,email" format:"email" example:"test@gmail.com"`
	Password string `json:"password" validate:"required,min=8,max=72" minLength:"8" maxLength:"72" example:"Password123"`
}

// Что должно быть в ответе.
// Для каждого поля нужно заполнить примеры для документации.
type Response struct {
    // Пример заполнения.
    // format, example, ...: https://github.com/swaggo/swag/blob/master/README.md
	AccessToken  string    `json:"accessToken" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"`
	RefreshToken uuid.UUID `json:"refreshToken" example:"289abe45-5920-4366-a12a-875ddb422ace"`
}

// Handle -.
//
// @Summary		Краткое название вашего обработчика (желательно на английском).
// @Description	Полное описание обработчика.
// @Tags			Тэги. Можете ставить на своё усмотрение. Можно ставить названия используемых сервисов (без приписки Service): auth/orders/users.
// @Param			request	body	Request	true	"Описание"
// @Accept			json
// @Produce		json
// @Success		200	{object}	Response
// @Failure		400	{object}	echo.HTTPError
// @Failure		404	{object}	echo.HTTPError
// @Failure		500	{object}	echo.HTTPError
// @Router			/example [post]
func (h *handler) Handle(c echo.Context, in Request) error {
    // Здесь вместо образца напишите свой код для обработки запроса.
	exampleData, err := h.exampleService.Login(c.Request().Context(), in.Email, in.Password)
	if err != nil {
        // Здесь обрабатываются ошибки, возвращаемые сервисом (обратите внимание на разные коды ответа).
		if errors.Is(err, example.ErrUserNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, err.Error())
		} else if errors.Is(err, example.ErrWrongPassword) {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

    // Обратите особое внимание на то, какой код ответа должен возвращать ваш обработчик.
	return c.JSON(http.StatusOK, Response{
		AccessToken:  authData.AccessToken,
		RefreshToken: authData.Session.RefreshToken,
	})
}
```

3. Посмотрите на комментарии в коде, не забудьте их стереть. Чеклист (убедитесь, что вы просмотрели следующие участки кода и заполнили их по инструкции):
   * `type handler struct {...}`
   * `type Request struct {...}`
   * `type Response struct {...}`
   * `func (h *handler) Handle(c echo.Context, in Request) error {...}`
   * документация (комментарии) перед `func (h *handler) Handle`
   * название пакета (первая строчка в файле), импорты
4. Интеграция обработчика в приложение
   * Добавьте в [`internal/app/handler.go`](../internal/app/handler.go) фабричный метод для создания вашего обработчика. Пример:
```go
func (app *App) PostExampleHandler() api.Handler {
    if app.postExampleHandler != nil {
        return app.postExampleHandler
    }

    // Предполагается, что у вас есть метод для получения нужного сервиса, например, app.ExampleService()
    app.postExampleHandler = post_example.New(app.ExampleService())
    return app.postExampleHandler
}
```
   * Зарегистрируйте маршрут в [`internal/app/router.go`](../internal/app/router.go). Добавьте новый маршрут в функцию `configureRouter` (если нужно, добавьте проверку прав).
   * Добавьте ручку как поле в структуру `App` в файле [`internal/app/app.go`](../internal/app/app.go).
```go
    postExampleHandler  api.Handler
```

## Организация сервисного слоя и репозиториев
Каждая функциональность (например, аутентификация, заказы, пользователи) реализуется в собственном пакете:
* [`internal/service`](../internal/service): содержит бизнес-логику.
* [`internal/repo`](../internal/repo): содержит доступ к данным и работу с БД.

При написании ручки обращайтесь к сервису, а не к репозиторию напрямую.

## Интеграция новых сервисов и репозиториев
Добавьте новые сервисы и репозитории как поля структуры `App` в файле [`internal/app/app.go`](../internal/app/app.go).
Не забудьте добавить для них конструкторы в файлы [`internal/app/db.go`](../internal/app/db.go) (для репозиториев), [`internal/app/service.go`](../internal/app/service.go) (для сервисов).

## Дополнительные ссылки
* Документация по валидатору: [ссылка](https://github.com/go-playground/validator/blob/master/README.md).
* Документация по заполнению Swagger: [ссылка](https://github.com/swaggo/swag/blob/master/README.md).