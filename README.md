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
- [Go:v1.23.2](https://github.com/golang/go)
- [PostgreSQL:v17.4](https://github.com/postgres/postgres)
- [Docker:v28.0.4](https://github.com/docker)
- [Docker compose:v2.34.0](https://github.com/docker/compose)

### Вспомогательные
- [PGX:v5.7.4(db driver)](https://github.com/jackc/pgx)
- [Goose:v3.24.2(migrations)](https://github.com/pressly/goose)
- [Squirrel:v1.5.4(query builder)](https://github.com/Masterminds/squirrel)

---

## Начало работы

**Перед началом работы сверьтесь с актуальными версиями технологий.**
**Также проверьте наличие конфигов, .env и make файла, а также файлы с миграциями**

### Быстрый старт

Выполните команду
```sh
make quick-up
```

Для завершения:
```sh
make quick-down
```

### Запуск локально с тестовыми бэкендами

Выполните команду
```sh
make lockal-up
```

Для завершения:
```sh
make lockal-down
```

### Запуск без тестовых бэкендов

Выполните команду
```sh
make no-b-up
```

Для завершения:
```sh
make no-b-down
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
