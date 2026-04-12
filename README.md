# TaskFlow

A task management REST API built in Go. Users can register, log in, create projects, and manage tasks inside those projects.

---

## 1. Overview

**Tech stack**

| | |
|---|---|
| Language | Go 1.24 |
| Router | chi v5 |
| Database | PostgreSQL 16 |
| Auth | JWT + bcrypt (cost 12) |
| Config | Viper |
| Logging | `log/slog` |
| Migrations | golang-migrate |

---

## 2. Architecture Decisions

The code is split into four layers:

- `domain/` — data models and repository interfaces. No external dependencies.
- `infrastructure/` — PostgreSQL, JWT, config, and logger implementations.
- `service/` — all business logic. Depends only on domain interfaces, not on the database or HTTP layer.
- `interfaces/http/` — chi handlers, middleware, and routing.

This separation keeps each layer focused and makes the service layer easy to test with in-memory mocks without needing a real database running.

**A few choices worth calling out**

- No ORM. SQL is written by hand. The schema is straightforward enough that an ORM would add more noise than value.
- Migrations run as a separate Docker Compose service (`migrate/migrate`) rather than being baked into the binary. This keeps the binary small and makes the migration step visible in compose logs.
- No refresh tokens. The assignment calls for a 24-hour JWT, so adding refresh token logic would be out of scope.

---

## 3. Running Locally

You only need Docker installed.

```bash
git clone https://github.com/pawannn/taskflow-pawan-kalyan
cd taskflow-pawan-kalyan
cp .env.example .env
docker compose up
```

The API will be available at **http://localhost:1337**.

When you run `docker compose up`:
1. PostgreSQL starts and passes a health check
2. The migrate service runs all migrations and seed data, then exits
3. The Go API server starts

```bash
# Rebuild after code changes
docker compose up --build

# Tear everything down (including the database volume)
docker compose down -v
```

---

## 4. Running Migrations

Migrations run automatically on startup. If you want to run them manually against a local database:

```bash
migrate -path db/migrations \
  -database "postgres://taskflow:taskflow@localhost:5432/taskflow?sslmode=disable" \
  up
```

To roll back:

```bash
migrate -path db/migrations \
  -database "postgres://taskflow:taskflow@localhost:5432/taskflow?sslmode=disable" \
  down
```

---

## 5. Test Credentials

The seed file (`db/migrations/000002_seed_data.up.sql`) creates three users. They all share the same password.

| Name | Email | Password |
|---|---|---|
| test | test@example.com | Test@Pass1! |
| pawan kalyan | pawan@gmail.com | Test@Pass1! |
| jhon | jhon@gmail.com | Test@Pass1! |

Quick login check:

```bash
curl -s -X POST http://localhost:1337/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"Test@Pass1!"}' | jq .
```

---

## 6. API Reference

A Postman collection covering all endpoints is included at `docs/taskflow.postman_collection.json`. Import it into Postman and set the `base_url` and `token` variables to get started quickly.

The full endpoint list is also documented below.

All protected endpoints require `Authorization: Bearer <token>`.

**Response format**
```json
{
  "req_id": "uuid",
  "status_code": 200,
  "client_message": "...",
  "data": {}
}
```

**Error format**
```json
{
  "req_id": "uuid",
  "status_code": 400,
  "error": "validation failed",
  "fields": { "email": "is required" }
}
```

**Status codes**

| Code | Meaning |
|---|---|
| 200 | OK |
| 201 | Created |
| 204 | No content (deletes) |
| 400 | Validation error — includes `fields` map |
| 401 | Missing or invalid token |
| 403 | Valid token, but you don't own this resource |
| 404 | Not found |
| 409 | Conflict (e.g. email already registered) |
| 500 | Internal server error |

---

### Auth

#### POST `/auth/register`
```json
// Request
{ "name": "pawan", "email": "pawan@gmail.com", "password": "Test@Pass1!" }

// 201
{ "data": { "id": "uuid", "name": "pawan", "email": "pawan@gmail.com", "created_at": "..." } }
```

#### POST `/auth/login`
```json
// Request
{ "email": "pawan@gmail.com", "password": "Test@Pass1!" }

// 200
{ "data": { "token": "<jwt>", "user": { "id": "...", "name": "...", "email": "..." } } }
```

---

### Projects

#### GET `/projects?page=1&limit=10`
Lists projects the logged-in user owns or has tasks assigned to them in.

#### POST `/projects`
```json
// Request
{ "name": "Website Redesign", "description": "Optional" }
// 201 — returns the created project
```

#### GET `/projects/:id?page=1&limit=10`
Returns the project along with its paginated tasks. Only accessible to the owner or an assignee.

#### PATCH `/projects/:id`
All fields are optional. Owner only.
```json
{ "name": "New Name", "description": "New description" }
// 200 — returns the updated project
```

#### DELETE `/projects/:id`
Owner only. Deletes the project and all its tasks. Returns `204`.

#### GET `/projects/:id/stats`
```json
// 200
{
  "data": {
    "total": 3,
    "status_counts": { "todo": 1, "in_progress": 1, "done": 1 },
    "assignee_counts": [
      { "assignee_id": "uuid", "count": 2 },
      { "assignee_id": null, "count": 1 }
    ]
  }
}
```

---

### Tasks

#### GET `/projects/:id/tasks?status=todo&assignee=uuid&page=1&limit=10`
Supports filtering by `status` (`todo`, `in_progress`, `done`) and `assignee` (UUID).

#### POST `/projects/:id/tasks`
Project owner only.
```json
// Request
{
  "title": "Design homepage",
  "description": "Optional",
  "priority": "high",
  "assignee_id": "uuid",
  "due_date": "2026-05-01"
}
// 201 — returns the created task
```

#### PATCH `/tasks/:id`
All fields optional. Project owner or task assignee only.
```json
{ "title": "Updated", "status": "done", "priority": "low", "due_date": "2026-06-01" }
// 200 — returns the updated task
```

#### DELETE `/tasks/:id`
Project owner or task creator only. Returns `204`.

---

## 7. Running Tests

13 integration tests covering auth, project, and task endpoints. Tests use in-memory mock repositories so no database is needed.

```bash
go test ./tests/integration/... -v
```

---

## 8. What I'd Do With More Time

- **Refresh tokens** — short-lived access tokens paired with long-lived refresh tokens stored in the database
- **`GET /projects/:id/members`** — useful for any UI that lets you pick an assignee from a list
- **Database-backed integration tests** using testcontainers — the current tests use mocks which are fast but won't catch SQL-level bugs
- **OpenAPI spec** — generate it from route definitions rather than maintaining docs by hand
- **CI pipeline** — lint, test, build, and push image on every PR
- **Soft deletes** — `deleted_at` on tasks and projects instead of hard deletes
- **Request-level timeout middleware** — individual DB calls have timeouts but the overall request currently does not
