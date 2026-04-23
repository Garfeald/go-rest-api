# REST API — Менеджер задач

REST API сервис для управления задачами, написанный на Go с использованием PostgreSQL.

## Технологии

- **Go** 1.26.2
- **PostgreSQL** 15 (через Docker)
- [sqlx](https://github.com/jmoiron/sqlx) — расширение стандартного `database/sql` с поддержкой маппинга в структуры
- [lib/pq](https://github.com/lib/pq) — драйвер PostgreSQL для Go
- **Docker Compose** — для запуска PostgreSQL и API

## Структура проекта

```
rest-api/
├── cmd/api/
│   └── main.go               # Точка входа приложения
├── internal/
│   ├── database/
│   │   ├── database.go        # Подключение к базе данных
│   │   └── tasks.go           # CRUD-операции над задачами
│   ├── handlers/
│   │   └── handlers.go        # HTTP-обработчики (handlers)
│   └── models/
│       └── task.go            # Модели данных (Task, CreateTaskInput, UpdateTaskInput)
├── sql/
│   └── init.sql               # SQL-скрипт инициализации БД (создание таблицы, сид-данные)
├── docker-compose.yml         # Конфигурация Docker Compose (PostgreSQL + API)
├── Dockerfile                 # Многоэтапная сборка Go-приложения
├── go.mod / go.sum            # Зависимости Go-модуля
└── README.md
```

## Модель данных

### Task

| Поле          | Тип         | Описание                   |
|---------------|-------------|----------------------------|
| `id`          | `int`       | Уникальный идентификатор   |
| `title`       | `string`    | Название задачи            |
| `description` | `string`    | Описание задачи            |
| `completed`   | `bool`      | Статус выполнения          |
| `created_at`  | `timestamp` | Дата создания              |
| `updated_at`  | `timestamp` | Дата последнего обновления |

### API-модели

- **CreateTaskInput** — `title`, `description`, `completed` (для создания задачи)
- **UpdateTaskInput** — `title`, `description`, `completed` (все поля указатели, для частичного обновления)

## API-эндпоинты

| Метод  | Маршрут          | Описание                     |
|--------|------------------|------------------------------|
| `GET`  | `/tasks`         | Получить все задачи          |
| `POST` | `/tasks/create`  | Создать новую задачу         |
| `GET`  | `/tasks/{id}`    | Получить задачу по ID        |
| `PUT`  | `/tasks/{id}`    | Обновить задачу по ID        |
| `DELETE`| `/tasks/{id}`   | Удалить задачу по ID         |

Все ответы возвращаются в формате JSON.

### Примеры запросов

**Создание задачи:**

```bash
curl -X POST http://localhost:8080/tasks/create \
  -H "Content-Type: application/json" \
  -d '{"title": "Новая задача", "description": "Описание", "completed": false}'
```

**Получение всех задач:**

```bash
curl http://localhost:8080/tasks
```

**Получение задачи по ID:**

```bash
curl http://localhost:8080/tasks/1
```

**Обновление задачи (частичное):**

```bash
curl -X PUT http://localhost:8080/tasks/1 \
  -H "Content-Type: application/json" \
  -d '{"completed": true}'
```

**Удаление задачи:**

```bash
curl -X DELETE http://localhost:8080/tasks/1
```

## Переменные окружения

| Переменная      | По умолчанию                                                              | Описание               |
|-----------------|---------------------------------------------------------------------------|------------------------|
| `DATABASE_URL`  | `postgres://taskuser:taskpass@localhost:5433/tasksdb?sslmode=disable`     | URL подключения к БД   |
| `SERVER_PORT`   | `8080`                                                                    | Порт HTTP-сервера      |

## Быстрый старт

### Предварительные требования

- [Go](https://go.dev/dl/) >= 1.26
- [Docker](https://www.docker.com/) + Docker Compose

### Вариант 1: Запуск через Docker Compose (рекомендуется)

Поднимает PostgreSQL и API вместе:

```bash
docker-compose up --build
```

API будет доступен на `http://localhost:8080`.

Для фонового запуска:

```bash
docker-compose up --build -d
```

### Вариант 2: Локальная разработка

#### 1. Запуск PostgreSQL

```bash
docker-compose up -d postgres
```

База данных будет доступна на `localhost:5432`:

| Параметр   | Значение    |
|------------|-------------|
| Host       | `localhost` |
| Port       | `5432`      |
| Database   | `tasksdb`   |
| User       | `taskuser`  |
| Password   | `taskpass`  |

> **Примечание:** Порт PostgreSQL в docker-compose — `5432`, а в Go-приложении по умолчанию указан `5433`. При локальной разработке задайте `DATABASE_URL` явно:

```bash
export DATABASE_URL="postgres://taskuser:taskpass@localhost:5432/tasksdb?sslmode=disable"
```

#### 2. Установка зависимостей

```bash
go mod download
```

#### 3. Запуск приложения

```bash
go run cmd/api/main.go
```

## Остановка

```bash
docker-compose down
```

Для удаления данных тома:

```bash
docker-compose down -v
```
