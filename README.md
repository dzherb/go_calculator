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

Из корня проекта выполните команду:

```shell
docker-compose up
```

После сборки образов и запуска контейнеров сервис будет доступен по адресу http://127.0.0.1:8080

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

#### Ответ (HTTP 201):
```json
{
  "id": 1
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
curl --location '127.0.0.1:8080/api/v1/calculate' \
--header 'Content-Type: application/json' \
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
  "error": "expected one of the methods: POST"
}
```
## Запуск тестов

Убедитесь, что находитесь в каталоге calculator_services. Затем запустите команду:
```shell
go test ./... 
```
