# TaskFlow

A task management REST API built in Go. Users register, log in, create projects, and manage tasks within those projects.

---

## 1. Overview

**Tech stack**

| | |
|---|---|
| Language | Go 1.26 |
| Router | chi v5 |
| Database | PostgreSQL 16 |
| Auth | JWT + bcrypt (cost 12) |
| Config | Viper |
| Logging | `log/slog` |
| Migrations | golang-migrate |

---

## 2. Architecture Decisions

The code is split into four layers:

- `domain/` — models and repository interfaces. No dependencies on anything external.
- `infrastructure/` — PostgreSQL, JWT, config, logger.
- `service/` — all business logic. Depends only on domain interfaces, not on database or HTTP.
- `interfaces/http/` — chi handlers, middleware, routing.

This separation means the service layer can be tested with in-memory mocks and no database running (which is exactly how the integration tests work).

**Choices worth noting**

- No ORM. SQL is written by hand. The schema is simple enough that an ORM adds noise, not value.
- Migrations run as a separate Docker Compose service (`migrate/migrate`) rather than being embedded in the binary. The binary stays small and dependency-free; the migration step is visible in compose logs.
- No refresh tokens. The assignment specifies a 24-hour JWT. Adding refresh logic would be scope creep.

---

## 3. Running Locally

You only need Docker.

```bash
git clone https://github.com/pawannn/taskflow-pawan-kalyan
cd taskflow-pawan-kalyan
cp .env.example .env
docker compose up
```

API is available at **http://localhost:1337**.

On `docker compose up`:
1. PostgreSQL starts and passes a health check
2. `migrate/migrate` runs all migrations including seed data, then exits
3. The Go API starts

```bash
# Rebuild after code changes
docker compose up --build

# Tear down completely
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

The seed file (`db/migrations/000002_seed_data.up.sql`) creates three users. All use the same password.

| Name | Email | Password |
|---|---|---|
| Test User | test@example.com | Test@Pass1! |
| Alice Johnson | alice@example.com | Test@Pass1! |
| Bob Smith | bob@example.com | Test@Pass1! |

Quick login:

```bash
curl -s -X POST http://localhost:1337/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"Test@Pass1!"}' | jq .
```

---

## 6. API Reference

All protected endpoints require `Authorization: Bearer <token>`.

**Response envelope**
```json
{
  "req_id": "uuid",
  "status_code": 200,
  "client_message": "...",
  "data": {}
}
```

**Error response**
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
{ "data": { "id": "uuid", "name": "pawan", "email": "pawan@example.com", "created_at": "..." } }
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
Lists projects the authenticated user owns or is assigned tasks in.

#### POST `/projects`
```json
// Request
{ "name": "Website Redesign", "description": "Optional" }
// 201 — returns created project
```

#### GET `/projects/:id?page=1&limit=10`
Returns the project and its paginated tasks. Must be owner or assignee to access.

#### PATCH `/projects/:id`
All fields optional. Owner only.
```json
{ "name": "New Name", "description": "New description" }
// 200 — returns updated project
```

#### DELETE `/projects/:id`
Owner only. Cascades to all tasks. Returns `204`.

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
Supports filtering by `status` (todo | in_progress | done) and `assignee` (UUID).

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
// 201 — returns created task
```

#### PATCH `/tasks/:id`
All fields optional. Project owner or task assignee only.
```json
{ "title": "Updated", "status": "done", "priority": "low", "due_date": "2026-06-01" }
// 200 — returns updated task
```

#### DELETE `/tasks/:id`
Project owner or task creator only. Returns `204`.

---

## 7. Running Tests

13 integration tests covering auth, project, and task endpoints. All tests use in-memory mock repositories — no database needed.

```bash
go test ./tests/integration/... -v
```

---

## 8. What I'd Do With More Time

- **Refresh tokens** — short-lived access tokens paired with long-lived refresh tokens stored in the database
- **`GET /projects/:id/members`** — needed for any real UI that lets you pick an assignee
- **Database-backed integration tests** using testcontainers — the current tests use mocks which are fast but don't catch SQL issues
- **OpenAPI spec** — generate it from route definitions rather than maintaining docs by hand
- **CI pipeline** — lint, test, build, push image
- **Soft deletes** — `deleted_at` on tasks and projects instead of hard deletes
- **Request-level timeout middleware** — individual DB calls have 10-second timeouts but the overall request has none
