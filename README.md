# Sport Assistance Backend

Backend-сервис для мобильного приложения и CRM Sport Assistance.

## Что это за проект
Проект реализует API для:
- регистрации и авторизации пользователей;
- работы с JWT access/refresh токенами;
- хранения пользователей, справочников и бизнес-данных в PostgreSQL;
- хранения активных access-токенов в Redis;
- применения SQL-миграций через `goose`.

## Текущий статус
Сейчас в проекте уже есть:
- базовая инфраструктура (`Go + Gin + Postgres + Redis + Docker`);
- auth-эндпоинты (`registration`, `login`, `refresh`, `logout`);
- миграции справочников и сущностей;
- роли и permissions в БД.

## Технологический стек
- Go 1.25
- Gin
- PostgreSQL (pgx)
- Redis
- Goose (миграции)
- Docker Compose

## Структура проекта
- `cmd/main.go` — точка входа.
- `internal/application` — сборка приложения и запуск сервера.
- `internal/handlers` — HTTP-обработчики.
- `internal/services` — бизнес-логика.
- `internal/repositories` — работа с БД.
- `internal/middlewares` — middleware (auth/cors).
- `migrations` — SQL-миграции.
- `docs` — OpenAPI/Swagger файлы.
- `pkg` — конфиг, логгер, DB/Redis коннекторы, server-обертка.

## Быстрый запуск (Docker)
1. Поднять контейнеры:
```bash
docker compose up -d --build
```
2. Проверить статус:
```bash
docker compose ps
```
3. Проверить логи приложения:
```bash
docker compose logs -f app
```

Что произойдет при запуске:
- поднимутся `postgres`, `redis`, `app`, `swagger`;
- `app` дождется healthcheck Postgres/Redis;
- внутри `app` автоматически выполнятся миграции (`goose up`).

## Порты и сервисы
| Контейнер | Сервис | Порт в контейнере | Порт на хосте | Назначение |
|---|---|---:|---:|---|
| `sport_assistance_app` | app | 8080 | 8080 | HTTP API |
| `fitness_postgres` | postgres | 5432 | 5432 | PostgreSQL |
| `redis-server` | redis | 6379 | 6379 | Redis |
| `sport_assistance_swagger` | swagger | 8080 | 8081 | Swagger UI |

## Точки доступа
- API: `http://localhost:8080`
- Healthcheck: `GET http://localhost:8080/ping`
- Swagger UI: `http://localhost:8081`
- OpenAPI YAML: `http://localhost:8081/docs/openapi.yaml`

## Переменные окружения
Основной шаблон: `.env.example`.

Ключевые группы:
- `SERVER_*` — порт и таймауты HTTP-сервера.
- `DB_*` — подключение к PostgreSQL.
- `GOOSE_*` — настройки миграций.
- `SECURITY_JWT_*` — секреты и TTL токенов.
- `REDIS_*` — подключение к Redis.
- `LOG_LEVEL`, `SWAGGER_ENABLED`.

## Запуск без Docker
1. Поднять отдельно Postgres и Redis.
2. Скопировать `.env.example` в `.env` и заполнить значения.
3. Применить миграции:
```bash
goose up
```
4. Запустить приложение:
```bash
go run ./cmd/main.go
```

## Миграции
Все миграции находятся в `migrations/`.

Текущее покрытие миграциями:
- базовые справочники;
- пользователи и связанные сущности;
- refresh tokens;
- список городов;
- роли/permissions.

Актуальные миграции по ролям и доступам:
- `migrations/00011_roles_permissions.sql`:
  - создает `roles`, `permissions`, `role_permissions`;
  - добавляет `users.role_id`.
- `migrations/00012_seed_roles_permissions.sql`:
  - заполняет роли (`guest`, `client`, `assistant`, `admin`);
  - заполняет список permissions;
  - заполняет связи `role_permissions`.

### Полезные команды goose
```bash
goose status
goose up
goose down
goose up-to 12
```

## Роли и permissions
В проекте используется RBAC-модель:
- роль хранится в `users.role_id`;
- права роли задаются в `role_permissions`.

Роли:
- `guest`
- `client`
- `assistant`
- `admin`

Примеры permissions:
- пользовательские: `profile.read.own`, `orders.read.own`, `matches.invite`;
- CRM: `crm.requests.read`, `crm.bookings.manage`;
- админские: `admin.users.manage`, `admin.finance.read`.

Важно:
- суффикс `.own` — это конвенция имени права (например, доступ только к своим данным);
- фактическая проверка ownership должна выполняться в бизнес-логике API.

## Текущие API-маршруты
Публичные (`/api/v1/auth`):
- `POST /registration`
- `POST /login`
- `POST /refresh`
- `POST /logout`

Служебный:
- `GET /ping`

## Диагностика и эксплуатация
Логи:
```bash
docker compose logs -f
docker compose logs -f app
```

Остановка:
```bash
docker compose down
```

Остановка с удалением данных БД/Redis:
```bash
docker compose down -v
```

Проверка таблиц в Postgres:
```bash
docker exec -it fitness_postgres psql -U fitness_user -d fitness_db -c "\dt"
```

## Известные нюансы
- На Windows следите за кодировкой файлов миграций: SQL-файлы должны быть UTF-8 **без BOM**.
- `go test ./...` может падать на `go vet`-предупреждениях в существующем коде; для проверки компиляции можно запускать:
```bash
go test -vet=off ./...
```

# Права

| Permission                   |      Гость     | Клиент |  Ассистент  | Администратор |
| ---------------------------- | :------------: | :----: | :---------: | :-----------: |
| **Аутентификация и профиль** |                |        |             |               |
| auth.register / login        |        ✅       |    ✅   |      ✅      |       ✅       |
| profile.view.own             |    ✅   |    ✅   |      ✅      |       ✅       |
| profile.edit.own             |    ⚠️ (огр.)   |    ✅   |      ❌      |       ❌       |
| profile.view.any             |        ❌       |    ❌   |      ✅      |       ✅       |
| profile.block / unblock      |        ❌       |    ❌   |      ❌      |       ✅       |
| change.user.subscription     |        ❌       |    ❌   |      ❌      |       ✅       |
| **Спортивный план**          |                |        |             |               |
| sport_plan.view              |   ⚠️ preview   |    ✅   |      ❌      |       ✅       |
| sport_plan.generate.ai       |        ❌       |    ✅   |      ❌      |       ❌       |
| sport_plan.edit              |        ❌       |    ❌   |      ❌      |       ❌       |
| progress.view                |        ❌       |    ✅   |      ❌      |       ✅       |
| **Матчи и расписание**       |                |        |             |               |
| schedule.view                |   ⚠️ частично  |    ✅   |      ❌      |       ✅       |
| activity.create              |        ❌       |    ✅   |      ❌      |       ❌       |
| match.create                 |        ❌       |    ✅   |      ❌      |       ❌       |
| match.invite.users           |        ❌       |    ✅   |      ❌      |       ❌       |
| match.confirm.participation  |        ✅       |    ✅   |      ❌      |       ❌       |
| match.enter.result           |        ❌       |    ✅   |      ❌      |       ❌       |
| match.manage.any             |        ❌       |    ❌   |      ✅      |       ✅       |
| **Рейтинг и друзья**         |                |        |             |               |
| rating.view                  |        ✅       |    ✅   |      ❌      |       ✅       |
| friends.view                 |  ⚠️ ограничено |    ✅   |      ❌      |       ❌       |
| friends.add                  | ⚠️ после матча |    ✅   |      ❌      |       ❌       |
| **Услуги и бронирования**    |                |        |             |               |
| booking.request              |        ❌       |    ✅   |      ❌      |       ❌       |
| booking.manage               |        ❌       |    ❌   |      ✅      |       ✅       |
| booking.cancel               |        ❌       |    ✅   |      ✅      |       ✅       |
| **Питание и экипировка**     |                |        |             |               |
| shop.view                    |   ⚠️ preview   |    ✅   |      ❌      |       ✅       |
| order.create                 |        ❌       |    ✅   |      ❌      |       ❌       |
| order.process                |        ❌       |    ❌   |      ✅      |       ❌       |
| order.view.own               |        ❌       |    ✅   |      ❌      |       ❌       |
| order.view.any               |        ❌       |    ❌   |      ✅      |       ✅       |
| **Кошелёк и платежи**        |                |        |             |               |
| wallet.view                  |        ❌       |    ✅   |      ❌      |       ✅       |
| wallet.topup                 |        ❌       |    ✅   |      ❌      |       ❌       |
| wallet.reserve               |        ❌       |    ✅   |      ❌      |       ❌       |
| wallet.writeoff              |        ❌       |    ❌   |      ❌      |       ❌       |
| wallet.force.adjust          |        ❌       |    ❌   |      ❌      |       ✅       |
| **Подписка**                 |                |        |             |               |
| subscription.view            |        ✅       |    ✅   |      ❌      |       ✅       |
| subscription.purchase        |        ✅       |    ✅   |      ❌      |       ❌       |
| subscription.cancel          |        ❌       |    ✅   |      ❌      |       ❌       |
| **Чат**                      |                |        |             |               |
| chat.view                    |   ⚠️ preview   |    ✅   |      ✅      |       ✅       |
| chat.send                    |        ❌       |    ✅   |      ✅      |       ✅       |
| chat.files.send              |        ❌       |    ✅   |      ✅      |       ✅       |
| **Опрос самочувствия**       |                |        |             |               |
| survey.fill                  |        ❌       |    ✅   |      ❌      |       ❌       |
| survey.view.own              |        ❌       |    ✅   |      ❌      |       ❌       |
| survey.view.any              |        ❌       |    ❌   |      ✅      |       ✅       |
| **CRM – справочники**        |                |        |             |               |
| specialists.view             |        ❌       |    ❌   |      ✅      |       ✅       |
| specialists.manage           |        ❌       |    ❌   | ⚠️ по праву |       ✅       |
| sport_objects.view           |        ❌       |    ❌   |      ✅      |       ✅       |
| sport_objects.manage         |        ❌       |    ❌   | ⚠️ по праву |       ✅       |
| services.manage              |        ❌       |    ❌   |      ❌      |       ✅       |
| **Новости**                  |                |        |             |               |
| news.view                    |        ✅       |    ✅   |      ❌      |       ✅       |
| news.manage                  |        ❌       |    ❌   | ⚠️ по праву |       ✅       |
| **Аналитика и финансы**      |                |        |             |               |
| analytics.view               |        ❌       |    ❌   |      ❌      |       ✅       |
| finance.view                 |        ❌       |    ❌   |      ❌      |       ✅       |
| finance.export               |        ❌       |    ❌   |      ❌      |       ✅       |
| **Системные функции**        |                |        |             |               |
| assistants.manage            |        ❌       |    ❌   |      ❌      |       ✅       |
| system.notifications.receive |        ❌       |    ❌   |      ❌      |       ✅       |
| admin.logs.view              |        ❌       |    ❌   |      ❌      |       ✅       |
