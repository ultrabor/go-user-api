# Go User API

REST API for user management written in Go with Gin. Layered structure: **delivery → services → storage**. Two interchangeable storages: PostgreSQL and in-memory.

## Features

- CRUD for users: create, get by ID, list, update, delete
- List with pagination (`page`, `limit`) and filters (`name`, `age`)
- Partial update: send only the fields you want to change
- Request logging middleware with `log/slog`
- Graceful shutdown on `SIGINT` / `SIGTERM`
- `404` for missing users, `400` for invalid input

## Quick start (in-memory)

```bash
git clone https://github.com/ultrabor/go-user-api.git
cd go-user-api
STORAGE=memory go run ./cmd
```

The server listens on `http://localhost:8080`.

## Run with PostgreSQL

PostgreSQL is the default storage.

1. Create the table:

   ```bash
   psql -h localhost -U postgres -f migration/001_create_users_table.sql
   ```

2. Set the connection (defaults in parentheses):

   | Variable | Default |
   |---|---|
   | `DB_HOST` | `localhost` |
   | `DB_PORT` | `5432` |
   | `DB_USER` | `postgres` |
   | `DB_PASSWORD` | `123` |
   | `DB_NAME` | `postgres` |
   | `DB_SSLMODE` | `disable` |

3. Run:

   ```bash
   go run ./cmd
   ```

## API

| Method | Path | Description | Success |
|---|---|---|---|
| `POST` | `/users` | Create a user | `201` `{"id": 1}` |
| `GET` | `/users` | List users | `200` |
| `GET` | `/users/{id}` | Get a user | `200` |
| `PUT` | `/users/{id}` | Update name and/or age | `200` |
| `DELETE` | `/users/{id}` | Delete a user | `204` |

Query parameters for `GET /users`: `page` (default `1`), `limit` (default `20`), `name`, `age`.

Errors are returned as `{"error": "message"}`.

## Examples

```bash
# create
curl -X POST http://localhost:8080/users \
  -H "Content-Type: application/json" \
  -d '{"name":"Ivan","age":25}'

# list with filters
curl "http://localhost:8080/users?name=Ivan&age=25&page=1&limit=20"

# get
curl http://localhost:8080/users/1

# update only age
curl -X PUT http://localhost:8080/users/1 \
  -H "Content-Type: application/json" \
  -d '{"age":26}'

# delete
curl -X DELETE http://localhost:8080/users/1
```

Ready-made requests for the VS Code REST Client / JetBrains HTTP Client are in [`examples/`](examples).

## Project structure

```
cmd/main.go                    entry point
internal/
  app/                         wiring, HTTP server, graceful shutdown
  config/                      environment variables
  delivery/                    Gin handlers, routes, logging middleware
  services/                    business logic
  domain/                      User model, UserStore interface, errors
  storage/
    memory/                    in-memory UserStore
    postgres/                  PostgreSQL UserStore (database/sql + lib/pq)
migration/                     SQL schema
examples/                      .http request files
```

`services` depends only on the `domain.UserStore` interface, so the storage is chosen at startup in `internal/app` with the `STORAGE` variable.
