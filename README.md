# Сервис для текстовых фрилансеров


## Предварительная проверка заданий

<a href=" ./../../../actions/workflows/1_helloworld.yml" >![1. Согласована и сформулирована тема курсовой]( ./../../actions/workflows/1_helloworld.yml/badge.svg)</a>

<a href=" ./../../../actions/workflows/2_usecase.yml" >![2. Usecase]( ./../../actions/workflows/2_usecase.yml/badge.svg)</a>

<a href=" ./../../../actions/workflows/3_data_model.yml" >![3. Модель данных]( ./../../actions/workflows/3_data_model.yml/badge.svg)</a>

<a href=" ./../../../actions/workflows/4_prototype_store_and_view.yml" >![4. Прототип хранение и представление]( ./../../actions/workflows/4_prototype_store_and_view.yml/badge.svg)</a>

<a href=" ./../../../actions/workflows/5_prototype_analysis.yml" >![5. Прототип анализ]( ./../../actions/workflows/5_prototype_analysis.yml/badge.svg)</a> 

<a href=" ./../../../actions/workflows/6_report.yml" >![6. Пояснительная записка]( ./../../actions/workflows/6_report.yml/badge.svg)</a>

<a href=" ./../../../actions/workflows/7_app_is_ready.yml" >![7. App is ready]( ./../../actions/workflows/7_app_is_ready.yml/badge.svg)</a>


## Данные для входа

Регистрация и логин полностью рабочие, так что можно создавать своих юзеров, при этом есть заготовлен тестовый сет с заказами и откликами + админом, который загружается в БД при запуске, если коллекция юзеров пустая.

Пароль для всех тестовых юзеров: `password123`
Далее будут указываться только почты для входа

### Админ

`admin@mail.com`

### Заказчики

- `client1@mail.com`
- `client2@mail.com`

### Фрилансеры

- `freelancer1@mail.com`
- `freelancer2@mail.com`
- `freelancer3@mail.com`
- `freelancer4@mail.com`


## Инструкция для запуска
1. Создать файл `.env`, можно скопировать [`.env.example`](.env.example) (файл уже создан, этот шаг можно пропустить).
2. Выполнить команду `make compose-up` ИЛИ `docker-compose up -d`.
3. Открыть в браузере страницу `127.0.0.1:1025`.