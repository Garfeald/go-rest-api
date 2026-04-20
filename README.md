# REST API — Менеджер задач

REST API сервис для управления задачами, написанный на Go с использованием PostgreSQL.

## Технологии

- **Go** 1.26.2
- **PostgreSQL** 15 (через Docker)
- [sqlx](https://github.com/jmoiron/sqlx) — расширение стандартного `database/sql` с поддержкой маппинга в структуры
- [lib/pq](https://github.com/lib/pq) — драйвер PostgreSQL для Go
- **Docker Compose** — для запуска PostgreSQL

## Структура проекта

```
rest-api/
├── cmd/                        # Точка входа приложения
├── internal/
│   ├── database/
│   │   ├── database.go         # Подключение к базе данных
│   │   └── tasks.go            # CRUD-операции над задачами
│   └── models/
│       └── task.go             # Модели данных (Task, CreateTaskInput, UpdateTaskInput)
├── sql/
│   └── init.sql                # SQL-скрипт инициализации БД (создание таблицы, сид-данные)
├── docker-compose.yml          # Конфигурация Docker Compose для PostgreSQL
├── go.mod / go.sum             # Зависимости Go-модуля
└── README.md
```

## Модель данных

### Task

| Поле         | Тип        | Описание                     |
|--------------|------------|------------------------------|
| `id`         | `int`      | Уникальный идентификатор     |
| `title`      | `string`   | Название задачи              |
| `description`| `string`   | Описание задачи              |
| `completed`  | `bool`     | Статус выполнения            |
| `created_at` | `timestamp`| Дата создания                |
| `updated_at` | `timestamp`| Дата последнего обновления   |

### API-модели

- **CreateTaskInput** — `title`, `description`, `completed` (для создания задачи)
- **UpdateTaskInput** — `title`, `description`, `completed` (все поля указатели, для частичного обновления)

## Быстрый старт

### Предварительные требования

- [Go](https://go.dev/dl/) >= 1.26
- [Docker](https://www.docker.com/) + Docker Compose

### 1. Запуск PostgreSQL

```bash
docker-compose up -d
```

База данных будет доступна на `localhost:5432`:

| Параметр   | Значение    |
|------------|-------------|
| Host       | `localhost` |
| Port       | `5432`      |
| Database   | `tasksdb`   |
| User       | `taskuser`  |
| Password   | `taskpass`  |

Строка подключения:

```
postgres://taskuser:taskpass@localhost:5432/tasksdb?sslmode=disable
```

При первом запуске автоматически выполняется `sql/init.sql` — создаются таблица `tasks` и три тестовые записи.

### 2. Установка зависимостей

```bash
go mod download
```

### 3. Запуск приложения

```bash
go run cmd/main.go
```

## Остановка

```bash
docker-compose down
```

Для удаления данных тома:

```bash
docker-compose down -v
```
