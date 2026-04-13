# TaskFlow

## 1. Overview

TaskFlow is a task management REST API that allows users to register, authenticate, create projects, and manage tasks within those projects. It supports assigning tasks, tracking their status and priority, and organizing work efficiently across multiple projects.

**Tech stack used:**

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

The system is designed with a clean, layered architecture to ensure separation of concerns, testability, and maintainability.

## 2. Architecture Decisions

The project follows a layered architecture to keep responsibilities clearly separated and the codebase maintainable as it grows.

- `domain/` — Contains core models and repository interfaces. This layer has no external dependencies, making it stable and easy to reason about.
- `infrastructure/` — Implements external concerns such as PostgreSQL, JWT, configuration, and logging. This layer satisfies the interfaces defined in the domain.
- `service/` — Holds all business logic. It depends only on the domain layer, not on HTTP or database implementations. This makes it easy to unit test using mocks.
- `interfaces/http/` — Handles HTTP transport (routing, middleware, request/response handling). It translates HTTP requests into service calls and formats responses.

**Why this structure?**
- Keeps business logic independent of frameworks and databases
- Makes testing easier (mock interfaces instead of real DB)
- Improves readability and scalability as features grow
- Avoids tight coupling between layers

**Key decisions and tradeoffs:**

##### **1. No ORM (raw SQL instead):**
- Chose explicit SQL for better control and performance.
- **Tradeoff**: More boilerplate and manual query management.

##### **2. JWT without refresh tokens:**
- Implemented simple 24-hour access tokens as the assignment calls.
- **Tradeoff**: Less flexible session management, but simpler and within scope.

##### **3. Migrations via separate service:**
- Runs as part of Docker Compose instead of being embedded in the app.
- **Tradeoff**: Slightly more setup complexity, but clearer lifecycle and smaller app binary.

##### **4. Layered architecture over simplicity:**
- Even though this is a small project, structured layering was used.
- **Tradeoff**: Slight initial overhead, but better long-term maintainability.

---

## 3. Running Locally

**You only need Docker installed.**

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
| test | test@example.com | password123 |
| pawan kalyan | pawan@gmail.com | password123 |
| jhon | jhon@gmail.com | password123 |

Quick login check:

```bash
curl -s -X POST http://localhost:1337/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"password123"}' | jq .
```

---

## 6. API Reference

A Postman collection covering all endpoints is available here: 

- **Collection**: [docs/taskflow.postman_collection.json](docs/taskflow.postman_collection.json).
- **Environment**: [docs/taskflow.postman_collection.json](docs/taskflow.postman_environment.json).

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
| 400 | Validation error - includes `fields` map |
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
{ "name": "pawan", "email": "pawan@gmail.com", "password": "password123" }

// 201
{ "data": { "id": "uuid", "name": "pawan", "email": "pawan@gmail.com", "created_at": "..." } }
```

#### POST `/auth/login`
```json
// Request
{ "email": "pawan@gmail.com", "password": "password123" }

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

## 7. What I'd Do With More Time

This implementation focuses on delivering a complete, clean, and working system within the given time constraints. Some decisions were simplified intentionally to prioritize correctness and clarity.

### Shortcuts taken

- **Simplified authentication**
  Only access tokens are implemented (no refresh tokens), as per assignment scope.

- **Mock-based integration tests**
  Tests use in-memory repositories instead of a real database. This keeps them fast but does not validate actual SQL queries.

---

### Improvements and additions

- **Refresh token mechanism**:
  Introduce short-lived access tokens with long-lived refresh tokens for better session management.

- **Database-backed integration tests**:
  Use tools like Testcontainers to run tests against a real PostgreSQL instance and catch query-level issues.

- **Project members & roles (RBAC)**:
  Add a proper members system with roles (owner, contributor) instead of relying only on ownership and task assignment.

- **OpenAPI / Swagger documentation**:
  Auto-generate API documentation to improve developer experience and reduce manual maintenance.

- **Observability**
  Add metrics (Prometheus) and tracing (OpenTelemetry) for better monitoring and debugging.

- **Soft deletes**
  Replace hard deletes with `deleted_at` fields for recoverability and auditing.
---

Overall, the focus was on building a solid, maintainable foundation. With more time, the system could be extended toward production-grade robustness and scalability.
