# sport-assistance

## Назначение
Backend-сервис для управления спортивными активностями. Запуск проекта выполняется через Docker Compose.

1. Запустите все контейнеры:
```bash
docker compose up -d --build
```
3. Проверьте состояние сервисов:
```bash
docker compose ps
```

## Порты контейнеров
| Контейнер | Сервис | Порт в контейнере | Порт на хосте | Назначение |
|---|---|---:|---:|---|
| `sport_assistance_app` | `app` | `8080` | `8080` | HTTP API |
| `fitness_postgres` | `postgres` | `5432` | `5432` | PostgreSQL |
| `redis-server` | `redis` | `6379` | `6379` | Redis |
| `sport_assistance_swagger` | `swagger` | `8080` | `8081` | Swagger UI |

## Точки доступа
- API: `http://localhost:8080`
- Swagger UI: `http://localhost:8081`
- OpenAPI файл: `http://localhost:8081/docs/openapi.yaml`

## Основные команды
Запуск:
```bash
docker compose up -d
```

Запуск с пересборкой:
```bash
docker compose up -d --build
```

Просмотр логов:
```bash
docker compose logs -f
```

Остановка контейнеров:
```bash
docker compose down
```

Остановка и удаление томов данных:
```bash
docker compose down -v
```

## Примечания
- Сервис `app` ожидает доступность `postgres` и `redis` и запускается после их успешного healthcheck.
- Для локального запуска вне Docker потребуется корректно настроенный `.env`.
