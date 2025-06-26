# Balancer project
Balancer project - это балансировщик нагрузки, который принимает входящие HTTP-запросы и распределяет их по пулу бэкенд-серверов.

---

## Содержание
- [Технологии](#технологии)
- [Начало работы](#начало-работы)
- [Архитектура и реализация](#архитектура-и-реализация)
- [API](#api)
- [Тестирование](#тестирование)
- [To do](#to-do)

---

## Технологии
### Основные
- [Go](https://github.com/golang/go)
- [PostgreSQL](https://github.com/postgres/postgres)
- [Docker](https://github.com/docker)
- [Docker compose](https://github.com/docker/compose)
- [Graylog](https://graylog.org/)

### Вспомогательные
- [PGX](https://github.com/jackc/pgx)
- [Goose](https://github.com/pressly/goose)
- [Squirrel](https://github.com/Masterminds/squirrel)

---

## Начало работы

### .env
```txt
PG_PORT=1234
PG_DATABASE_NAME=your_database_name
PG_USER=your_pg_user
PG_PASSWORD=your_awesome_password1234
MIGRATION_DIR=./your_migrations/migrations

GRAYLOG_PASSWORD_SECRET=somepasswordsalt
GRAYLOG_ROOT_PASSWORD_SHA2=your_sha256_password_hash
```

### config
```yaml
server:                                 # balancer server settings
  host: localhost                       # balancer server host
  port: 8080                            # balancer server port

ticker_rate_sec: 1                      # ticker rate for backends heath check

default_limits:                         # defoult settings for clients limits
  capacity: 5                           # defoult token capacity
  rate_per_sec: 0.1                     # defoult token refiil rate

backend_list:                           # backends settings
  - backend_url: http://localhost:8081  # backend url
    config:                             # current backend settings
      health:                           # current backend health check settings
        method: HEAD                    # health check method
        url: /                          # health check url
```

### Полный запуск проекта

Выполните команды
```sh
make fullsetup-up
make migration-up
```

Для завершения:
```sh
make fullsetup-down
```

### Запуск локально с тестовыми бэкендам

Выполните команду
```sh
make local-up
make migration-up
go run cmd/main.go -config-path config/local/config.yaml
```

Для завершения:
```sh
make local-down
```

---

## Архитектура и реализация

```cmd``` - корень всего проекта, инициализация конфига, репозитория, сервисов, api.

```internal/api/handler``` - слой HTTP endpoit'ы для взаимодействия

```internal/config``` - парсинг и логика инициализации конфига

```internal/integration-suite``` - итеграционные тесты

```internal/model``` - entities

```internal/repository``` - взаимодействие с базой данных

```internal/service``` - usecase, бизнес логика

```internal/service/in-memory-cache``` - реализация кеша и логика работы с ним

```internal/service/interfaces``` - интерфейсы всех сервисов

```internal/service/limits-manager``` - логика работы с изменением rate лимита клиентов

```internal/service/strategy``` - логика алгоритмов для распределения нагрузки и health checks бэкендов

```internal/service/token-manager``` - модуль для ограничения частоты запросов (rate-limiting) на основе алгоритма Token Bucket

---

## API

API подразделяется на два типа:

**API для отправки запросов на backend-серверы через балансировщик**
**API для управления rate лимитами пользователей**

---

### 1. API отправки запросов на backend-серверы

Эндпоинт, через который балансировщик проксирует входящие HTTP-запросы к backend-серверам:
`/` — проксирование запросов на один из backend-серверов согласно стратегии балансировки

### 2. API управления rate лимитом пользователей

| Метод  | URL                | Назначение                         | Пример запроса (JSON)                                                         | Пример ответа (JSON)                                                      | Код ответа      |
|--------|--------------------|------------------------------------|-------------------------------------------------------------------------------|---------------------------------------------------------------------------|-----------------|
| POST   | `/limits/create`   | Создать новый лимит                | ```{ "client_id": "123.5.6.7:890", "capacity": 10, "rate_per_sec": 1 }```     | —                                                                         | `201 Created`   |
| GET    | `/limits/get`      | Получить лимит по client_id        | ```{ "client_id": "123.5.6.7:890" }```                                        | ```{ "client_id": "123.5.6.7:890", "capacity": 10, "rate_per_sec": 1 }``` | `200 OK`        |
| PUT    | `/limits/update`   | Обновить лимит                     | ```{ "client_id": "123.5.6.7:890", "capacity": 10, "rate_per_sec": 1 }```     | —                                                                         | `204 NoContent` |
| DELETE | `/limits/delete`   | Удалить лимит                      | ```{ "client_id": "123.5.6.7:890" }```                                        | —                                                                         | `204 NoContent` |

---

## Тестирование
Для тестирования используется простой echo сервер ```test-backend```
Он принимает json файлы и выводит метод, uri и отправленный json
Request:
```sh
{
    "exaple":123
}
```
Reasponse:
```sh
{
    "server_address": "loclahost:8081",
    "http_method": "GET",
    "request_uri": "/",
    "your_message": {
        "exaple": 123
    }
}
```

Реализовано unit-тестирование

Интеграционные тесты:
internal/integration-suite

---

## To do
- [x] Алгоритм распределения list connections
- [x] Unit тесты
- [x] Интеграционное тестирование
- [X] Логгирование
- [ ] Трейсинг
- [ ] Метрики
- [ ] Добавить LRU cache
