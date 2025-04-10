# Работа с авторизацией в проекте

## Общие сведения

Используется механизм `accessToken` + `refreshToken`. Суть в том, что `accessToken` короткоживущий и было бы неудобно каждые 30 минут заново логиниться через ручку `/auth/login` (здесь и дальше рассматриваются пути на бэке, с фронта к ним еще прибавляется префикс `api`), поэтому для предотвращения этого есть `refreshToken`, с помощью которого по ручке `/auth/refresh` можно обновить пару токенов незаметно от юзера. Кроме того, это позволяет быть залогиненым на разных устройствах, в разных браузерах и тд именно за счет того, что под сессию долгоживущую (в нашем случае это месяц) выделяется отдельный `refreshToken`.

Механизм работы нашей системы примерно такой:
- `/auth/login` позволяет залогиниться по данным пользователя, которые передаются в теле POST запроса. Возращается `accessToken` и `refreshToken` (рефреш еще проставляется в http-only cookie, чем мы и будем пользоваться, а возврат в теле нужен для того, если вдруг будет мобильное приложение, где понятия нормальных кук нет). `accessToken` живет 30 минут и в нем закодированы данные о `userID` и `systemRole` и именно он подставляется в заголовок `Authorization` в виде `Bearer <тут сам токен>` всех будущих HTTP-запросов на "защищенные" авторизацией ручки. `accessToken` фронту можно сохранить у себя локально, в то время как `refreshToken` летает по кукам и причем только на `/auth` роуты и его надо максимально не палить. `refreshToken` живет 30 дней, но как его использовать - далее
- возникает вопрос, "а что делать по истечению 30 минут?". Ответ таков: посмотреть, что `accessToken` действительно протух, отослать `refreshToken` на `/auth/refresh` и получить новую пару токенов, причем отправленный в запросе `refreshToken` станет невалидным и надо будет использовать тот, что придет в ответе. То есть по сути потребность логиниться заново может потребоваться только если 30 дней не заходить на сайт, потому что протухнут уже оба токена, в остальных же случаях можно получить новую пару незаметно для клиента (этим будет заниматься код на фронте)
- на будущее также предусмотрен ендпоинт `/auth/logout`, который позволяет дропнуть сессию по `refreshToken`, с помощью этого в случае утечки рефреша можно будет подропать все сессии и дать залогиниться только реальному юзеру

## Необходимые знания по коду

В [internal/app/router.go](../internal/app/router.go) в `configureRouter` может пригодиться навешивать middleware (дальше mw, это такая штука, которая отрабатывает до запроса и чаще всего требуется много, где, поэтому ее вот так выделяют. она может решить не допускать до ручки, как в случае с отсутствием авторизации, может что-то сделать с контекстом и тд). Пример как это сделано для ручки `GET /admin`

```
adminGroup := handler.Group("/admin", app.AuthMW().UserIdentity())
{
	adminGroup.GET("", app.GetAdminHandler().Handle, app.AuthMW().AdminRole())
}
```

Здесь на всю группу ручек накидывается mw `UserIdentity`, который парсит `accessToken` и записывает в `c echo.Context` информацию о роли и информацию об ID юзера. Благодаря этому в ручке можно через `c.Get(<нужный ключ>)` получить эти значения и сразу закастить к нужному типу, потому что на это есть гарантия. Отдельно же для `GET /admin` навешивается мидлвейр, проверяющий, что юзер админ `AdminRole`. Здесь три важных замечания до самого примера с кодом из ручки `GET /admin` [internal/api/get_admin/handler.go](../internal/api/get_admin/handler.go):
- `AdminRole` рассчитывает, что данные о роли уже в контексте, поэтому следите, чтобы `UserIdentity` в цепочке был раньше. Здесь мы на всю группу навесили сначала `UserIdentity`, а потом уже на отдельную ручку `AdminRole`, поэтому все окей
- `<нужный ключ>` это константы, которые надо записывать в пакете `mw`
- чтобы линтер не ругался на кастинг к нужному типу без проверки, нужно прописать комментарий `//nolint:forcetypeassert`

```
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
```

Это довольно примитивная ручка, которая вытаскивает из контекста запроса `c echo.Context` проставленные UserIdentity мидлвейром параметры и возвращает их в ответе. Ну и до самой ручки допускаются только админы.
Ниже пример, как в сваггере, который после запуска доступен по [http://localhost/api/swagger/index.html](http://localhost/api/swagger/index.html), был получен ответ из ручки `GET /admin`. Предварительно была проведена авторизация через `auth/login` и `accessToken` был записан в авторизационное поле в формате `Bearer <accessToken>`, скриншот поля тоже приложил
![image](https://github.com/user-attachments/assets/5f67bc2b-59a8-41bd-829e-2f13fcb08053)
![swag](https://github.com/user-attachments/assets/35195b9f-3af2-417d-9aa1-b5e606fa4ba8)

