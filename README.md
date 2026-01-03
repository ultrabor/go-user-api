# Go User API

Простой и расширяемый REST API на Go для управления пользователями.  
Проект построен по модульной (чистой) архитектуре: **handlers → services → storage**.  
Поддерживаются **in-memory** и **PostgreSQL** хранилища, логирование через middleware.

---

## Возможности

- Создание пользователя  
- Получение пользователя по ID  
- Получение списка пользователей  
- Обновление пользователя  
- Удаление пользователя  
- Логирование HTTP-запросов  
- Хранилища:
  - In-memory
  - PostgreSQL

---

## Быстрый старт

Запуск с in-memory хранилищем:

```bash
git clone https://github.com/ultrabor/go-user-api.git
cd go-user-api
go run ./cmd
```

Сервер будет доступен по адресу:  
`http://localhost:8080`

---

## Запуск с PostgreSQL

### 1. Миграция базы данных

```bash
psql postgres://user:pass@localhost:5432/dbname \
  -f migration/001_create_users_table.sql
```

### 2. Переменные окружения

```bash
export STORAGE=postgres
export DATABASE_URL="postgres://user:pass@localhost:5432/dbname?sslmode=disable"
```

### 3. Запуск

```bash
go run ./cmd
```

---

## API Endpoints

| Метод  | Путь           | Описание |
|------:|----------------|----------|
| POST   | `/create`      | Создать пользователя |
| GET    | `/get/{id}`    | Получить пользователя |
| GET    | `/users`       | Получить всех пользователей |
| PUT    | `/update/{id}` | Обновить пользователя |
| DELETE | `/delete/{id}` | Удалить пользователя |

---

## Примеры запросов

### Создание пользователя

```bash
curl -X POST http://localhost:8080/create \
  -H "Content-Type: application/json" \
  -d '{"name":"Ivan","age":25}'
```

### Получение пользователя

```bash
curl http://localhost:8080/get/1
```

### Получение списка пользователей

```bash
curl http://localhost:8080/users
```

### Обновление пользователя

```bash
curl -X PUT http://localhost:8080/update/1 \
  -H "Content-Type: application/json" \
  -d '{"name":"Ivan Petrov","age":26}'
```

### Удаление пользователя

```bash
curl -X DELETE http://localhost:8080/delete/1
```

---

## Структура проекта

```
.
├── cmd
│   └── main.go
├── examples
│   ├── create.http
│   ├── delete.http
│   ├── get.http
│   ├── getAll.http
│   └── update.http
├── internal
│   ├── app
│   ├── config
│   ├── handlers
│   ├── middleware
│   ├── models
│   ├── server
│   ├── services
│   └── storage
│       ├── memory
│       └── postgres
├── migration
│   └── 001_create_users_table.sql
└── README.md
```

---

## Архитектура

- **handlers** — HTTP-слой  
- **services** — бизнес-логика  
- **storage** — интерфейс `UserStore` и его реализации  
- **middleware** — логирование (`log/slog`)  
- **models** — структуры данных

