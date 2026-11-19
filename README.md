# Task Management Application

Fullstack task management system built with Go (Gin), MySQL, Redis, and React Native (Expo) with TypeScript.

## Tech Stack

### Backend

- Go 1.25+
- Gin Web Framework
- MySQL 8.0 (`database/sql` + `go-sql-driver/mysql`)
- Redis 7.0 (`go-redis/v9`)
- Testing: `testify`, `go-sqlmock`

### Frontend

- React Native 0.88 / Expo SDK 58 (Expo Router)
- TypeScript
- Testing: Jest, `@testing-library/react-native`

---

## Project Structure

```text
task-management/
├── docker-compose.yml              # MySQL and Redis services
├── README.md
├── backend/
│   ├── cmd/server/main.go          # Application entrypoint & HTTP server
│   ├── internal/
│   │   ├── cache/                  # Redis caching implementation & key generation
│   │   ├── handler/                # HTTP request handlers & response helpers
│   │   ├── model/                  # Task structs, DTOs, and filter models
│   │   ├── repository/             # MySQL database queries & soft delete logic
│   │   └── service/                # Business logic & cache coordination
│   ├── migrations/
│   │   └── 001_create_tasks.sql    # DDL schema definition & seed data
│   ├── go.mod
│   └── go.sum
└── mobile/
    ├── src/
    │   ├── app/
    │   │   └── index.tsx           # Task management main screen
    │   ├── components/             # Reusable UI components (search, filter, card, modal, pagination)
    │   │   └── __tests__/          # Component tests
    │   └── services/api.ts         # Backend API client
    ├── package.json
    └── tsconfig.json
```

---

## Setup and Installation

### Prerequisites

- Go 1.25 or higher
- Node.js (v20+) and npm
- Docker and Docker Compose

### 1. Database and Cache (Docker)

Start MySQL and Redis using Docker Compose:

```bash
docker compose up -d
```

Containers started:

- MySQL: `localhost:3306` (database: `taskdb`, credentials: `root:root`)
- Redis: `localhost:6379`

The migration file in `backend/migrations/001_create_tasks.sql` is automatically executed on startup.

### 2. Run Backend

```bash
cd backend
go run cmd/server/main.go
```

The server runs on `http://localhost:8080`.

Optional environment variables:

- `PORT` (default: `8080`)
- `DB_HOST` (default: `localhost`)
- `DB_PORT` (default: `3306`)
- `DB_USER` (default: `root`)
- `DB_PASSWORD` (default: `root`)
- `DB_NAME` (default: `taskdb`)
- `REDIS_ADDR` (default: `localhost:6379`)
- `REDIS_PASSWORD` (default: empty)

### 3. Run Mobile / Web Frontend

In a separate terminal:

```bash
cd mobile
npm run web
```

To run on mobile platforms:

- Android: `npm run android`
- iOS: `npm run ios`

---

## Database Schema & Migrations

The database migration is located in `backend/migrations/001_create_tasks.sql`:

```sql
CREATE TABLE IF NOT EXISTS tasks (
  id INT AUTO_INCREMENT PRIMARY KEY,
  title VARCHAR(255) NOT NULL,
  description TEXT,
  status ENUM('todo','in_progress','done') NOT NULL DEFAULT 'todo',
  assignee VARCHAR(100),
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP NULL,
  UNIQUE KEY uq_title (title),
  INDEX idx_status (status),
  INDEX idx_assignee (assignee)
);
```

Design details:

- Soft delete: `deleted_at` column is updated with a timestamp on delete. All queries filter with `WHERE deleted_at IS NULL`.
- Unique title: `UNIQUE KEY uq_title (title)` prevents duplicate titles at the database layer.
- Indexes: `idx_status` and `idx_assignee` index frequent filter targets.

---

## API Reference

Base URL: `http://localhost:8080/api`

### Endpoints

| Method   | Endpoint         | Description                                         |
| -------- | ---------------- | --------------------------------------------------- |
| `GET`    | `/health`        | Health check                                        |
| `GET`    | `/api/tasks`     | Get paginated task list with filters (Redis cached) |
| `POST`   | `/api/tasks`     | Create new task (invalidates cache)                 |
| `GET`    | `/api/tasks/:id` | Get task by ID                                      |
| `PUT`    | `/api/tasks/:id` | Update task by ID (invalidates cache)               |
| `DELETE` | `/api/tasks/:id` | Soft delete task by ID (invalidates cache)          |

### GET /api/tasks Parameters

- `status` (string, optional): `todo`, `in_progress`, or `done`. Returns 400 if invalid.
- `keyword` (string, optional): Substring search in `title` or `description`.
- `assignee` (string, optional): Exact filter by assignee.
- `page` (int, default: `1`): Current page.
- `limit` (int, default: `10`, max: `100`): Items per page.
- `sort` (string, default: `created_at DESC`): Sorting field and direction (e.g., `created_at ASC`, `title ASC`). Whitelisted to prevent SQL injection.

### Response Format

Success (list):

```json
{
  "success": true,
  "data": [
    {
      "id": 1,
      "title": "Setup Docker environment",
      "description": "Initialize MySQL and Redis containers",
      "status": "done",
      "assignee": "devops",
      "created_at": "2026-11-19T09:00:00Z",
      "updated_at": "2026-11-19T09:00:00Z"
    }
  ],
  "meta": {
    "page": 1,
    "limit": 10,
    "total": 1,
    "total_pages": 1
  }
}
```

Error response:

```json
{
  "success": false,
  "error": {
    "code": "CONFLICT",
    "message": "A task with this title already exists"
  }
}
```

HTTP status codes used:

- `200 OK`: Successful query, update, or delete.
- `201 Created`: Resource successfully created.
- `400 Bad Request`: Validation failure or invalid parameters.
- `404 Not Found`: Task not found or already soft deleted.
- `409 Conflict`: Duplicate title error (MySQL 1062).
- `500 Internal Server Error`: Unexpected database or system error.

---

## Redis Caching Strategy

1. TTL: Responses from `GET /api/tasks` are cached for 60 seconds.
2. Cache Key: Normalized query parameters are appended to the key prefix:
   `tasks:status=<status>:keyword=<keyword>:assignee=<assignee>:page=<page>:limit=<limit>:sort=<sort>`
3. Invalidation: Any create, update, or delete operation scans and purges all keys matching `tasks:*`.
4. Fallback: If Redis is unreachable, the system logs a warning and queries MySQL directly.

---

## Bug Fixes and Improvements

1. Duplicate Title (HTTP 409):
   - MySQL error 1062 (`ER_DUP_ENTRY`) was previously bubbling up as a generic 500 error.
   - The repository detects duplicate key errors and maps them to `ErrDuplicateTitle`.
   - The handler returns HTTP 409 Conflict with a clear error payload.
   - The frontend modal captures 409 responses and renders an inline warning banner.

2. Refresh List After Update:
   - Creating, editing, or deleting a task in the UI immediately triggers a list refetch to keep data in sync.

3. Soft Delete:
   - `DELETE /api/tasks/:id` updates `deleted_at = CURRENT_TIMESTAMP`.
   - All `SELECT` queries enforce `WHERE deleted_at IS NULL`.
   - Deleted items are omitted from search results, pagination counts, and direct lookups.

---

## Testing

### Backend Unit Tests

Covers repository queries with `go-sqlmock`, service cache hit/miss/invalidation, and handler error codes:

```bash
cd backend
go test -v ./...
```

Tests included:

- `TestMySQLTaskRepository_Update_Success`
- `TestMySQLTaskRepository_Update_DuplicateTitle_ReturnsConflict`
- `TestMySQLTaskRepository_List_SearchAndHideSoftDeleted`
- `TestMySQLTaskRepository_SoftDelete`
- `TestTaskService_Search_CacheMissAndHit`
- `TestTaskService_Update_InvalidatesCache`
- `TestTaskService_CreateAndDelete_InvalidatesCache`
- `TestHandler_CreateTask_DuplicateTitle_Returns409`
- `TestHandler_UpdateTask_DuplicateTitle_Returns409`
- `TestHandler_UpdateTask_NotFound_Returns404`
- `TestHandler_GetTasks_InvalidStatus_Returns400`

### Frontend Component Tests

Covers search input, status filters, task card actions, and pagination navigation:

```bash
cd mobile
npm test
```

Tests included in `mobile/src/components/__tests__/task-components.test.tsx`:

- `TaskSearchBar`: Input changes and clear button handling.
- `TaskStatusFilter`: Selection and status change callbacks.
- `TaskCard`: Rendering details, assignee, status badge, edit and delete events.
- `TaskPagination`: Pagination state, button disabling, and page navigation.
