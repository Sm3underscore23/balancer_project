# Balancer project
Balancer project - это простой балансировщик нагрузки, который принимает входящие HTTP-запросы и распределяет их по пулу бэкенд-серверов.

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
**Также проверьте наличие конфигов, .env и make файла**

### Запуск проекта с помощью docker
Установите docker на свое устройство:
[Гайд по установке](https://docs.docker.com/engine/install/)

**Установите сопутствующие инстансы**
```sh
docker pull postgres:17.4
```
```sh
docker pull s3m23/balancer:v0.1
```
```sh
docker pull s3m23/backend-test:v0.1
```

Выполните команды для запуска контейнеров
```sh
make fullsetup-up
```

Установите goose и выполните миграции
```sh
make install-goose
```
```sh
make migration-up
```


Для завершения:
```sh
make fullsetup-down
```

---

## Архитектура и реализация

```cmd``` - корень всего проекта, инициализация конфига, репозитория, сервисов, api, происходит тут

```internal/api/handler``` - слой в котором находятся HTTP endpoit'ы для взаимодействия с проектом

```internal/config``` - парсинг и логика инициализации конфига

```internal/model``` - внутренние абстракции для взаимодействия с проектом и для его работы

```internal/repository``` - взаимодействие проекта с базой данных

```internal/service``` - бизнес логика проекта

Стоит пройтись по подробнее по реализованным сервисам:

```internal/service/in-memory-cache``` - реализация простого кеша и логики работы с ним

```internal/service/interfaces``` - интерфейсы всех сервисов

```internal/service/limits-manager``` - реализация логики работы с изменением rate лимита клиентов

```internal/service/pool-service``` - реализация pool'а бэкендов, health checks, а также реализация стратегий для распределения нагрузки на бэкенды

```internal/service/token-manager``` - модуль для ограничения частоты запросов (rate-limiting) на основе алгоритма Token Bucket

---

## API

API подразделяется на два типа:

**API для отправки запросов на backend-серверы через балансировщик**
**API для управления rate лимитами пользователей**

---

### 1. API отправки запросов на backend-серверы

Эндпоинт, через который балансировщик проксирует входящие HTTP-запросы к backend-серверам:

- `/` — проксирование запросов на один из backend-серверов согласно стратегии балансировки

### 2. API управления rate лимитом пользователей

| Метод  | URL                | Назначение                          | Пример запроса (JSON)                                                                 | Пример ответа (JSON)                                                            | Код ответа |
|--------|--------------------|-------------------------------------|---------------------------------------------------------------------------------------|----------------------------------------------------------------------------------|------------|
| POST   | `/limits/create`   | Создать новый лимит                | ```json<br>{ "client_id": "123.5.6.7:890", "capacity": 10, "rate_per_sec": 1 }```     | —                                                                                | `201 Created` |
| GET    | `/limits/get`      | Получить лимит по client_id        | ```json<br>{ "client_id": "123.5.6.7:890" }```                                        | ```json<br>{ "client_id": "123.5.6.7:890", "capacity": 10, "rate_per_sec": 1 }``` | `200 OK` |
| PUT    | `/limits/update`   | Обновить лимит                     | ```json<br>{ "client_id": "123.5.6.7:890", "capacity": 10, "rate_per_sec": 1 }```     | —                                                                                | `204 NoContent` |
| DELETE | `/limits/delete`   | Удалить лимит                      | ```json<br>{ "client_id": "123.5.6.7:890" }```                                        | —                                                                                | `204 NoContent` |

---

## Тестирование
Для тестирования я разработал простой echo сервер ```test-backend```
Он принимает json файлы и выводит метод и uri и отправленный json
Request:
```sh
{
    "exaple":123
}
```
Reasponse:
```sh
{
    "http_method": "GET",
    "request_uri": "/",
    "your_message": {
        "exaple": 123
    }
}
```
---

## To do
- [x] Добавить алгоритм распределения least connections
- [x] Добавить unit тесты
- [x] Добавить интеграционное тестирование
- [ ] Добавить метрики
- [ ] Добавить LRU cache
