# TaskFlow

A minimal task management API built in Go — users can register, log in, create projects, and manage tasks within those projects.

---

## 1. Overview

TaskFlow is a REST API implementing authentication, project management, and task tracking with role-based access control.

**Tech stack:**
- **Language:** Go 1.21+
- **Router:** chi v5
- **Database:** PostgreSQL 16
- **Auth:** JWT (golang-jwt/jwt v5) + bcrypt
- **Config:** Viper (reads `.env` or environment variables)
- **Logging:** `log/slog` (structured, per-concern log files + stdout)
- **Migrations:** golang-migrate CLI (`migrate/migrate` Docker image)

---

## 2. Architecture Decisions

**Layered / hexagonal structure.** The project is split into four layers:
- `domain/` — models and repository interfaces (no dependencies on infrastructure)
- `infrastructure/` — PostgreSQL implementations, JWT, config, logger
- `service/` — business logic; depends only on domain interfaces
- `interfaces/http/` — chi handlers, middleware, routing; depends only on service types

This makes the business logic testable without a database (see integration tests, which inject mock repositories).

**What I left out and why:**
- No ORM — SQL is written by hand for full control and visibility; the schema is simple enough that an ORM adds no value here.
- No refresh tokens — the assignment specifies a 24-hour JWT; adding refresh logic would exceed scope.
- No rate limiting — out of scope for this exercise but would be the first thing added before production exposure.
- No WebSocket/SSE — backend-only submission.

**Tradeoffs made:**
- `logger.New()` always writes to `logs/` files on disk. In production this would be replaced with a structured sink (e.g. Loki, CloudWatch). For this exercise it's acceptable.
- The `migrate/migrate` Docker image is used as a separate compose service rather than embedding the runner in the binary. This keeps the binary dependency-free and makes the migration step explicit and observable.

---

## 3. Running Locally

Requires: Docker and Docker Compose (nothing else).

```bash
git clone https://github.com/pawannn/taskflow-pawan-kalyan
cd taskflow-pawan-kalyan
cp .env.example .env
docker compose up
```

The API will be available at **http://localhost:1337**.

What happens on `docker compose up`:
1. PostgreSQL starts and passes its health check
2. `migrate/migrate` runs all migrations (schema + seed data) and exits
3. The Go API binary starts

To rebuild the API image after code changes:
```bash
docker compose up --build
```

To stop and remove containers + volumes:
```bash
docker compose down -v
```

---

## 4. Running Migrations

Migrations run automatically on `docker compose up` via the `migrate` service.

To run migrations manually against a local PostgreSQL instance:
```bash
# Install the CLI: https://github.com/golang-migrate/migrate/releases
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

The seed migration (`000002_seed_data.up.sql`) creates the following users:

| Name | Email | Password |
|---|---|---|
| Test User | `test@example.com` | `Test@Pass1!` |
| Alice Johnson | `alice@example.com` | `Test@Pass1!` |
| Bob Smith | `bob@example.com` | `Test@Pass1!` |

After `docker compose up`, log in with:
```bash
curl -s -X POST http://localhost:1337/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"Test@Pass1!"}' | jq .
```

---

## 6. API Reference

All protected endpoints require `Authorization: Bearer <token>`.

All responses follow this envelope:
```json
{
  "req_id": "uuid",
  "status_code": 200,
  "client_message": "...",
  "data": { ... }
}
```
Error responses:
```json
{
  "req_id": "uuid",
  "status_code": 400,
  "error": "validation failed",
  "fields": { "email": "is required" }
}
```

### Auth

#### POST `/auth/register`
```json
// Request
{ "name": "Jane Doe", "email": "jane@example.com", "password": "Test@Pass1!" }

// 201 Created
{ "data": { "id": "uuid", "name": "Jane Doe", "email": "jane@example.com", "created_at": "..." } }
```

#### POST `/auth/login`
```json
// Request
{ "email": "jane@example.com", "password": "Test@Pass1!" }

// 200 OK
{ "data": { "token": "<jwt>", "user": { "id": "...", "name": "...", "email": "..." } } }
```

### Projects

#### GET `/projects?page=1&limit=10`
Returns projects the authenticated user owns or has tasks in.

#### POST `/projects`
```json
{ "name": "My Project", "description": "Optional" }
// 201 Created — returns project object
```

#### GET `/projects/:id?page=1&limit=10`
Returns project details + its tasks. Requires the caller to be owner or assignee.

#### PATCH `/projects/:id`
```json
{ "name": "Updated Name", "description": "Updated" }
// 200 OK — owner only
```

#### DELETE `/projects/:id`
```
204 No Content — owner only, cascades to all tasks
```

#### GET `/projects/:id/stats`
Returns task counts broken down by status and by assignee.
```json
// 200 OK
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

### Tasks

#### GET `/projects/:id/tasks?status=todo&assignee=uuid&page=1&limit=10`
Filters: `status` (todo | in_progress | done), `assignee` (UUID).

#### POST `/projects/:id/tasks`
```json
{
  "title": "Design homepage",
  "description": "Optional",
  "priority": "high",
  "assignee_id": "uuid or null",
  "due_date": "2026-05-01"
}
// 201 Created
```

#### PATCH `/tasks/:id`
All fields optional. Caller must be project owner or task assignee.
```json
{ "title": "Updated", "status": "done", "priority": "low", "due_date": "2026-06-01" }
// 200 OK
```

#### DELETE `/tasks/:id`
Project owner or task assignee only. Returns `200 OK`.

### HTTP Status Codes

| Code | Meaning |
|---|---|
| 200 | Success |
| 201 | Created |
| 400 | Validation error (includes `fields` map) |
| 401 | Missing or invalid JWT |
| 403 | Valid JWT but insufficient permission |
| 404 | Resource not found |
| 409 | Conflict (e.g. email already registered) |
| 500 | Internal server error |

---

## 7. Running Tests

Integration tests cover the auth endpoints using an in-memory mock repository (no database required):

```bash
go test ./tests/integration/... -v
```

---

## 8. What I'd Do With More Time

**Shortcuts taken:**
- The `logs/` directory is always written to disk. A proper production setup would use structured log sinks (Loki, CloudWatch) configured per environment.
- `DELETE /tasks/:id` returns `200` with a body instead of `204 No Content` — the assignment says 204 but the existing error/response helpers assume a body. Easy to fix.
- No request-level timeout middleware — individual repository calls have 10-second timeouts but the overall request doesn't.

**What I'd add:**
- Refresh token flow (short-lived access tokens + long-lived refresh tokens stored in DB)
- Rate limiting on auth endpoints (login brute-force protection)
- `GET /users` or `GET /projects/:id/members` to support assignee lookup in a real UI
- Database-backed integration tests using testcontainers to cover the full stack
- OpenAPI/Swagger spec generated from route definitions
- Proper CI pipeline (lint, test, build, Docker push)
- Soft deletes on tasks/projects with `deleted_at` column
