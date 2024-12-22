# GO-калькулятор

## Описание
GO-калькулятор — это веб-сервис для вычисления результата числового выражения, переданного пользователем.

Поддерживаемые операции:
- Сложение (`+`)
- Вычитание (`-`)
- Умножение (`*`)
- Деление (`/`)
- Приоритет операций с использованием скобок

## Сборка и запуск

По умолчанию сервер запускается на *127.0.0.1:8080*.  
При необходимости адрес и порт можно изменить с помощью переменных окружения:

```shell
export GO_CALC_ADDR=0.0.0.0  # Установка адреса на 0.0.0.0
export GO_CALC_PORT=8081     # Установка порта на 8081
```

Для запуска приложения выполните следующую команду:

```shell
go run cmd/main.go
```

## Взаимодействие с сервисом

### Успешный запрос

#### Пример
```shell
curl --location '127.0.0.1:8080/api/v1/calculate' \
--header 'Content-Type: application/json' \
--data '{
  "expression": "(2+2)*2/2.5"
}'
```

#### Ответ (HTTP 200):
```json
{
  "result": 3.2
}
```

### Ошибка: невалидный JSON

#### Пример
```shell
curl --location '127.0.0.1:8080/api/v1/calculate' \
--header 'Content-Type: application/json' \
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
curl --location '127.0.0.1:8080/api/v1/calculate' \
--header 'Content-Type: application/json' \
--data '{
  "expression": "abc - 3"
}'
```

#### Ответ (HTTP 422):
```json
{
  "error": "expression is not valid"
}
```

### Ошибка: не передан параметр "expression" или он пустой

#### Пример
```shell
curl --location '127.0.0.1:8080/api/v1/calculate' \
--header 'Content-Type: application/json' \
--data '{
  "exp": "4 + 2"
}'
```

#### Ответ (HTTP 400):
```json
{
  "error": "expression is required"
}
```

### Ошибка: неверный HTTP-метод

#### Пример
```shell
curl --location '127.0.0.1:8080/api/v1/calculate' \
-X GET
--header 'Content-Type: application/json' \
--data '{
  "expression": "2 + 2"
}'
```

#### Ответ (HTTP 405):
```json
{
  "error": "expected POST method"
}
```
## Запуск тестов

Тесты можно запустить следующей командой:
```shell
go test ./... 
```
