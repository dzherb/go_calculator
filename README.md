# GO-калькулятор

## Описание
**GO-калькулятор** — это веб-сервис для распределенного вычисления результата числового выражения. Почти как обычный калькулятор, но с gRPC, PostgreSQL и аутентификацией, разве не круто?

Приложение развернуто по адресу: https://go-calc.dzherb.ru/

## Сборка и запуск

Из корня проекта выполните команду:

```shell
docker compose up
```

После сборки образов и запуска контейнеров приложение станет доступно по адресу http://127.0.0.1:8081

## Взаимодействие с сервисом

## Через интерфейс

После запуска docker compose можно сразу начать взаимодействие с сервисом прямо из браузера. Результат вычисления, статус обработки и ошибки отображаются в реальном времени: 

![Скриншот интерфейса клиента](readme_assets/ui_screenshot.png)

Предварительно понадобится пройти регистрацию:

![Скриншот страницы регистрации](readme_assets/ui_register_screenshot.png)

## Через curl

### Несколько слов об аутентификации

В приложении она работает на [jwt](https://jwt.io/introduction)-токенах. Время жизни токена определяется переменной окружения `ACCESS_TOKEN_TTL`, по-умолчанию равной одному часу.

После получения токена его необходимо прикреплять к заголовкам в виде:
`Authorization: Bearer ваш_токен`

Заголовок, очевидно, не требуется для путей `api/v1/auth/login` и `api/v1/auth/register`, во всех остальных случаях он ожидается.
___

### Регистрация

```shell
curl --location '127.0.0.1:8081/api/v1/auth/register' \
--header 'Content-Type: application/json' \
--data '{
  "username": "batman", "password": "I_hate_superman"
}'
```

#### Ответ (HTTP 200)

```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NDY5NjEyMjguNzg4Nzc2OSwiaWF0IjoxNzQ2OTU3NjI4Ljc4ODc3NjksInN1YiI6IjIifQ.IsC12YSkpxSJ23Sp0l9rewpL09oLhXGKzd3448OITXU",
  "user": {
    "id": 2,
    "username": "batman",
    "created_at": "2025-05-11T10:00:28.758033Z",
    "updated_at": "2025-05-11T10:00:28.758033Z"
  }
}
```


### Вход / получение нового токена
```shell
curl --location '127.0.0.1:8081/api/v1/auth/login' \
--header 'Content-Type: application/json' \
--data '{
  "username": "batman", "password": "I_hate_superman"
}'
```

#### Ответ (HTTP 200)

```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NDY5NjI2NzkuNDgwMjM5NCwiaWF0IjoxNzQ2OTU5MDc5LjQ4MDIzOTQsInN1YiI6IjIifQ.w6Hx9Jj2j4umUJDusie5VB1V7SLUPOFzW7RGP8vsVno",
  "user": {
    "id": 2,
    "username": "batman",
    "created_at": "2025-05-11T10:00:28.758033Z",
    "updated_at": "2025-05-11T10:00:28.758033Z"
  }
}
```

___

### Успешный запрос

#### Пример
```shell
curl --location '127.0.0.1:8081/api/v1/calculate' \
--header 'Content-Type: application/json' \
--header 'Authorization: Bearer ваш_токен' \
--data '{
  "expression": "(2+2)*2/2.5"
}'
```

#### Ответ (HTTP 201):
```json
{
  "id": 1
}
```

### Ошибка: не предоставлен access token
```shell
curl --location '127.0.0.1:8081/api/v1/calculate' \
--header 'Content-Type: application/json' \
--data '{
  "expression": "(2+2)*2/2.5"
}'
```

#### Ответ (HTTP 401):
```json
{
  "error": "token must be provided"
}
```

### Ошибка: невалидный JSON

#### Пример
```shell
curl --location '127.0.0.1:8081/api/v1/calculate' \
--header 'Content-Type: application/json' \
--header 'Authorization: Bearer ваш_токен' \
--data '{
  "exp
}'
```

#### Ответ (HTTP 500):
```json
{
  "error": "failed to parse request body: invalid character '\\n' in string literal"
}
```

### Ошибка: невалидное выражение

#### Пример
```shell
curl --location '127.0.0.1:8081/api/v1/calculate' \
--header 'Content-Type: application/json' \
--header 'Authorization: Bearer ваш_токен' \
--data '{
  "expression": "3 + a"
}'
```

#### Ответ (HTTP 422):
```json
{
  "error": "expression contains invalid token at position 5: а"
}
```

### Ошибка: не передан параметр "expression" или он пустой

#### Пример
```shell
curl --location '127.0.0.1:8081/api/v1/calculate' \
--header 'Content-Type: application/json' \
--header 'Authorization: Bearer ваш_токен' \
--data '{
  "exp": "4 + 2"
}'
```

#### Ответ (HTTP 400):
```json
{
  "error": "expression is empty"
}
```

### Ошибка: неверный HTTP-метод

#### Пример
```shell
curl --location '127.0.0.1:8081/api/v1/calculate' \
-X GET
--header 'Content-Type: application/json' \
--header 'Authorization: Bearer ваш_токен' \
--data '{
  "expression": "2 + 2"
}'
```

#### Ответ (HTTP 405):
```json
{
  "error": "expected one of the methods: POST"
}
```
## Настройка констант

Оркестратор берет значения констант из переменных окружения. Можно задать их отдельно перед запуском сервиса, но вариант проще — отредактировать значения прямо в **docker-compose.yml**

```yaml
TIME_ADDITION_MS: 100
TIME_SUBTRACTION_MS: 100
TIME_MULTIPLICATIONS_MS: 100
TIME_DIVISIONS_MS: 100
TASK_MAX_PROCESS_TIME_IN_MS: 30000
```

Первые четыре отвечают за время выполнения базовых операций, а последняя — за общее максимальное время обработки агентом задачи, после чего оркестратор отказывается принимать ее результат.

Также при необходимости можно поменять порты бэкенд-сервиса и клиента, это все задается в том же **docker-compose.yml**

## Схема взаимодействия сервисов

![Схема сервисов](readme_assets/services_schema.png)

## Запуск тестов

Убедитесь, что находитесь в каталоге calculator. Затем запустите команду:
```shell
make test
```
Будут выполнены как unit-тесты, так и интеграционные. Для последних автоматически поднимется временный Docker-контейнер с PostgreSQL.

## Схема БД

Ознакомиться со схемой базы данных можно [здесь](database.md).

Также доступна команда:
```bash
make schema
```

Убедитесь сперва, что установлены утилиты [migrate](https://github.com/golang-migrate/migrate) и [npm](https://www.npmjs.com/). После запуска будет обновлено содержание [database.md](database.md).

Для этого скрипт поднимет временный контейнер с PostgreSQL, `migrate` применит актуальные миграции (находятся в `internal/storage/migrations`), а затем `pg-mermaid` на основе актуальной схемы БД сгенерирует ER-диаграмму.