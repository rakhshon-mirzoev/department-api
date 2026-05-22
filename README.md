# Department API

API для управления подразделениями и сотрудниками компании.

## Технологии

- Go, net/http
- PostgreSQL, GORM
- Миграции через goose
- Docker, docker-compose

## Запуск

```bash
docker-compose up --build
```

API запустится на `http://localhost:8080`

## Эндпоинты

```
POST   /departments/                        создать подразделение
GET    /departments/{id}                    получить подразделение (depth, include_employees)
PATCH  /departments/{id}                    обновить название или родителя
DELETE /departments/{id}                    удалить (mode=cascade или mode=reassign)
POST   /departments/{id}/employees/         создать сотрудника
```
