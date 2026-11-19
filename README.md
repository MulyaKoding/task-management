# Fullstack Task Management Application

A production-ready Fullstack Task Management application built with **Go (Gin)**, **MySQL**, **Redis**, and **React Native (Expo)** with **TypeScript**.

Developed as part of the Fullstack Engineer Assessment, fulfilling all functional requirements, caching strategies, bug fixes, automated test suites, and documentation.

---

## 📌 Table of Contents
1. [Tech Stack](#-tech-stack)
2. [Project Architecture](#-project-architecture)
3. [Requirement Completion Matrix](#-requirement-completion-matrix)
4. [Getting Started](#-getting-started)
   - [Prerequisites](#prerequisites)
   - [Step 1: Start MySQL & Redis (Docker)](#step-1-start-mysql--redis-docker)
   - [Step 2: Run Backend](#step-2-run-backend)
   - [Step 3: Run Frontend](#step-3-run-frontend)
5. [Database Migration](#-database-migration)
6. [API Documentation](#-api-documentation)
7. [Redis Caching Strategy](#-redis-caching-strategy)
8. [Bug Fixes Implemented](#-bug-fixes-implemented)
9. [Automated Testing](#-automated-testing)
   - [Backend Unit Tests](#backend-unit-tests)
   - [Frontend Component Tests](#frontend-component-tests)

---

## 🚀 Tech Stack

### Backend
- **Language**: Go 1.25+
- **Framework**: [Gin Web Framework](https://github.com/gin-gonic/gin)
- **Database**: MySQL 8.0 (`database/sql` + `go-sql-driver/mysql`)
- **Cache**: Redis 7.0 (`go-redis/v9`)
- **Testing**: `testify`, `go-sqlmock`

### Frontend
- **Framework**: React Native 0.88 with [Expo SDK 58](https://expo.dev/) (Expo Router)
- **Language**: TypeScript 6.0+
- **Styling**: React Native StyleSheet & Custom Theming (Light/Dark mode)
- **Testing**: Jest, `@testing-library/react-native`, `@react-native/jest-preset`

---

## 📂 Project Architecture

```
task-management/
├── docker-compose.yml              # MySQL & Redis container definitions
├── backend/
│   ├── cmd/
│   │   └── server/
│   │       └── main.go             # Server entrypoint, router setup, CORS
│   ├── internal/
│   │   ├── cache/
│   │   │   └── task_cache.go       # Redis cache client & key generator
│   │   ├── handler/
│   │   │   ├── task_handler.go     # Gin HTTP handlers
│   │   │   ├── task_handler_test.go# Handler unit tests (409, 404, 400)
│   │   │   └── response.go         # Standardized API response format
│   │   ├── model/
│   │   │   └── task.go             # Domain entities, DTOs, response models
│   │   ├── repository/
│   │   │   ├── task_repository.go  # SQL queries (soft delete, filters, sorting)
│   │   │   └── task_repository_test.go # Repository tests (sqlmock)
│   │   └── service/
│   │       ├── task_service.go     # Business logic & cache orchestration
│   │       └── task_service_test.go# Service tests (cache hit/miss/invalidation)
│   ├── migrations/
│   │   └── 001_create_tasks.sql    # DDL schema with soft delete & unique key
│   ├── go.mod
│   └── go.sum
├── mobile/
│   ├── src/
│   │   ├── app/
│   │   │   └── index.tsx           # Main Tasks screen (Search, Filter, List, Modals)
│   │   ├── components/
│   │   │   ├── task-card.tsx       # Task card item with badges & actions
│   │   │   ├── task-modal.tsx      # Create/Edit modal with conflict handling
│   │   │   ├── task-pagination.tsx # Pagination controls
│   │   │   ├── task-search-bar.tsx # Search input with clear button
│   │   │   ├── task-status-filter.tsx # Status pill buttons (All, Todo, In Progress, Done)
│   │   │   └── __tests__/
│   │   │       └── task-components.test.tsx # Component unit tests
│   │   └── services/
│   │       └── api.ts              # API client methods with typed responses
│   ├── package.json
│   └── tsconfig.json
└── README.md
```

---

## ✅ Requirement Completion Matrix

| Assessment Task | Requirement Details | Implementation Status |
| :--- | :--- | :---: |
| **Task 1: Backend (40%)** | • Filtering (`status`, `keyword`, `assignee`, `page`, `limit`, `sort`)<br>• Implement `PUT /api/tasks/{id}`<br>• Implement soft `DELETE /api/tasks/{id}`<br>• Return consistent error responses | **Completed** |
| **Task 2: Redis (15%)** | • Cache `GET /api/tasks` for 60 seconds<br>• Invalidate cache after create/update/delete<br>• Cache key includes query parameters | **Completed** |
| **Task 3: Frontend (25%)** | • Search input with clear trigger<br>• Status filter pills (`All`, `To Do`, `In Progress`, `Done`)<br>• Pagination controls (prev/next, total count, pages)<br>• Edit modal (title, description, status, assignee)<br>• Loading states & empty states | **Completed** |
| **Task 4: Bug Fixes (10%)** | • Duplicate title returns HTTP 409 Conflict instead of 500<br>• Refresh list after update/create/delete<br>• Hide soft-deleted tasks from all queries (`deleted_at IS NULL`) | **Completed** |
| **Task 5: Testing (10%)** | • Backend: Update, Search, Cache invalidation<br>• Frontend: Component tests for search, filter, card, and pagination | **Completed** |

---

## 🛠 Getting Started

### Prerequisites
- [Go](https://golang.org/) (version 1.25 or later)
- [Node.js](https://nodejs.org/) (v20+ or v22+) & npm
- [Docker & Docker Compose](https://www.docker.com/) (for running MySQL & Redis)

### Step 1: Start MySQL & Redis (Docker)
In the project root directory, run:
```bash
docker compose up -d
```
This spins up:
- **MySQL 8.0** on `localhost:3306` with database `taskdb`
- Automatically executes `./backend/migrations/001_create_tasks.sql`
- **Redis 7.0** on `localhost:6379`

To check container status:
```bash
docker compose ps
```

### Step 2: Run Backend
Navigate to the `backend` directory and start the server:
```bash
cd backend
go run cmd/server/main.go
```
The Gin server will start on `http://localhost:8080`.

#### Environment Variables (Optional):
You can customize the connection using environment variables:
| Variable | Default Value | Description |
| :--- | :--- | :--- |
| `PORT` | `8080` | Port for the HTTP server |
| `DB_HOST` | `localhost` | MySQL host |
| `DB_PORT` | `3306` | MySQL port |
| `DB_USER` | `root` | MySQL user |
| `DB_PASSWORD` | `root` | MySQL password |
| `DB_NAME` | `taskdb` | MySQL database name |
| `REDIS_ADDR` | `localhost:6379` | Redis host:port address |
| `REDIS_PASSWORD`| `""` | Redis password |

### Step 3: Run Frontend
Open a new terminal, navigate to `mobile`, and start Expo:
```bash
cd mobile
npm run web
# Or run on Android/iOS:
# npm run android
# npm run ios
```
Open `http://localhost:8081` in your browser.

---

## 🗄 Database Migration

The database schema is located at `backend/migrations/001_create_tasks.sql`:

```sql
CREATE TABLE IF NOT EXISTS tasks (
  id INT AUTO_INCREMENT PRIMARY KEY,
  title VARCHAR(255) NOT NULL,
  description TEXT,
  status ENUM('todo','in_progress','done') NOT NULL DEFAULT 'todo',
  assignee VARCHAR(100),
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP NULL,            -- soft delete
  UNIQUE KEY uq_title (title),          -- duplicate title -> 409 Conflict
  INDEX idx_status (status),
  INDEX idx_assignee (assignee)
);
```

### Key Highlights:
1. **Soft Delete**: `deleted_at TIMESTAMP NULL`. Deleting a task sets `deleted_at = CURRENT_TIMESTAMP`. All retrieval and list queries enforce `WHERE deleted_at IS NULL`.
2. **Duplicate Prevention**: `UNIQUE KEY uq_title (title)` prevents duplicate task titles. If violated, MySQL error `1062` is caught and returned as HTTP `409 Conflict`.
3. **Optimized Indexes**: `idx_status` and `idx_assignee` provide fast lookups and high filtering performance.

---

## 📡 API Documentation

### Base URL: `http://localhost:8080/api`

### Standard Response Format:

#### Success Response:
```json
{
  "success": true,
  "data": { ... },
  "meta": {
    "page": 1,
    "limit": 10,
    "total": 5,
    "total_pages": 1
  }
}
```

#### Consistent Error Response:
```json
{
  "success": false,
  "error": {
    "code": "CONFLICT",
    "message": "A task with this title already exists"
  }
}
```

---

### Endpoints Overview:

| Method | Endpoint | Description |
| :--- | :--- | :--- |
| `GET` | `/health` | Health check endpoint |
| `GET` | `/api/tasks` | Get paginated & filtered list of tasks (Cached in Redis) |
| `POST` | `/api/tasks` | Create new task (Invalidates Redis cache) |
| `GET` | `/api/tasks/:id` | Get single task by ID |
| `PUT` | `/api/tasks/:id` | Update task details (Invalidates Redis cache) |
| `DELETE` | `/api/tasks/:id` | Soft delete task (Invalidates Redis cache) |

---

### Query Parameters for `GET /api/tasks`:
- `status` *(optional)*: Filter by status (`todo`, `in_progress`, `done`). Returns `400 Bad Request` if invalid.
- `keyword` *(optional)*: Search term across `title` and `description`.
- `assignee` *(optional)*: Filter by assignee name.
- `page` *(optional, default: 1)*: Current page number.
- `limit` *(optional, default: 10, max: 100)*: Items per page.
- `sort` *(optional, default: `created_at DESC`)*: Sort column & order (e.g., `created_at ASC`, `title ASC`). Whitelisted against SQL injection.

---

### HTTP Status Codes:
- `200 OK`: Request succeeded.
- `201 Created`: Resource successfully created.
- `400 Bad Request`: Validation failure (e.g. empty title, invalid status).
- `404 Not Found`: Task does not exist or has been soft-deleted.
- `409 Conflict`: Duplicate title conflict on create or update.
- `500 Internal Server Error`: Unexpected database/server error.

---

## ⚡ Redis Caching Strategy

1. **TTL (Time-To-Live)**:
   - `GET /api/tasks` responses are cached in Redis for **60 seconds**.
2. **Cache Key Pattern**:
   - The cache key includes all normalized query parameters:
     ```
     tasks:status=<status>:keyword=<keyword>:assignee=<assignee>:page=<page>:limit=<limit>:sort=<sort>
     ```
   - Example: `tasks:status=in_progress:keyword=auth:assignee=alice:page=1:limit=10:sort=created_at DESC`
3. **Cache Invalidation**:
   - Any write operation (`POST /api/tasks`, `PUT /api/tasks/:id`, `DELETE /api/tasks/:id`) automatically invalidates all task cache keys by scanning and deleting keys matching `tasks:*`.
4. **Resilience**:
   - If Redis is unavailable or disconnected, the service logs a warning and queries MySQL directly without failing the request.

---

## 🐞 Bug Fixes Implemented

1. **Duplicate Title Returns HTTP 409 Instead of 500**:
   - Previously, duplicate title insertions caused an unhandled internal server error (500).
   - MySQL error code `1062` (`ER_DUP_ENTRY`) is detected in the repository and mapped to `repository.ErrDuplicateTitle`.
   - The handler returns HTTP `409 Conflict` with `{ "success": false, "error": { "code": "CONFLICT", "message": "A task with this title already exists" } }`.
   - The frontend modal displays a clear conflict warning banner without crashing.

2. **Refresh List After Update**:
   - When a task is created, updated via the Edit Modal, or soft-deleted, the frontend immediately triggers `loadTasks()`, ensuring the UI list is always up-to-date with backend state.

3. **Hide Soft-Deleted Tasks**:
   - All `SELECT` queries enforce `WHERE deleted_at IS NULL`.
   - Soft-deleted items (`deleted_at IS NOT NULL`) are hidden from search results, pagination counts, and detail retrieval.

---

## 🧪 Automated Testing

### Backend Unit Tests
The backend test suite covers:
1. **Search & Filter** (`internal/repository/task_repository_test.go`):
   - Multi-field search, pagination, sort, and exclusion of soft-deleted records using `go-sqlmock`.
2. **Task Update & Duplicate Title Conflict** (`internal/repository/task_repository_test.go`):
   - Successful updates and MySQL error 1062 translation to `ErrDuplicateTitle`.
3. **Soft Delete** (`internal/repository/task_repository_test.go`):
   - Updating `deleted_at` timestamp and handling non-existent/already deleted records.
4. **Cache Miss, Cache Hit & Invalidation** (`internal/service/task_service_test.go`):
   - First search hits repo and sets cache.
   - Second search serves directly from cache without querying repo.
   - Task creation, update, and deletion trigger cache invalidation.
5. **HTTP Status Codes & Consistent Error Handling** (`internal/handler/task_handler_test.go`):
   - Verifies 409 Conflict, 404 Not Found, 400 Bad Request, and JSON response structure.

#### Run Backend Tests:
```bash
cd backend
go test -v ./...
```
**Result**:
```
PASS: TestHandler_CreateTask_DuplicateTitle_Returns409 (0.00s)
PASS: TestHandler_UpdateTask_DuplicateTitle_Returns409 (0.00s)
PASS: TestHandler_UpdateTask_NotFound_Returns404 (0.00s)
PASS: TestHandler_GetTasks_InvalidStatus_Returns400 (0.00s)
PASS: TestMySQLTaskRepository_Update_Success (0.00s)
PASS: TestMySQLTaskRepository_Update_DuplicateTitle_ReturnsConflict (0.00s)
PASS: TestMySQLTaskRepository_List_SearchAndHideSoftDeleted (0.00s)
PASS: TestMySQLTaskRepository_SoftDelete (0.00s)
PASS: TestTaskService_Search_CacheMissAndHit (0.00s)
PASS: TestTaskService_Update_InvalidatesCache (0.00s)
PASS: TestTaskService_CreateAndDelete_InvalidatesCache (0.00s)
ok  	task-app/internal/handler
ok  	task-app/internal/repository
ok  	task-app/internal/service
```

---

### Frontend Component Tests
The frontend test suite (`mobile/src/components/__tests__/task-components.test.tsx`) uses `@testing-library/react-native` and Jest to test:
1. **TaskSearchBar**:
   - Renders with custom placeholder.
   - Fires `onChangeText` callback when input text changes.
   - Shows clear button when value is non-empty and fires `onClear`.
2. **TaskStatusFilter**:
   - Renders all status options: "All", "To Do", "In Progress", "Done".
   - Handles selection and fires `onSelectStatus` with the chosen status.
3. **TaskCard**:
   - Renders title, description, status badge, and assignee.
   - Triggers `onEdit` when "Edit" button is pressed.
   - Triggers `onDelete` when "Delete" button is pressed.
4. **TaskPagination**:
   - Renders current page number and total task count.
   - Disables "Previous" on page 1.
   - Triggers `onPageChange` with next page on "Next" press.

#### Run Frontend Tests:
```bash
cd mobile
npm test
```
**Result**:
```
PASS src/components/__tests__/task-components.test.tsx
  Task Management UI Component Tests
    TaskSearchBar
      √ renders with placeholder and triggers onChangeText
      √ renders clear button when value is present and triggers onClear
    TaskStatusFilter
      √ renders all status options and calls onSelectStatus on selection
    TaskCard
      √ displays task details, assignee, status badge and triggers edit & delete
    TaskPagination
      √ renders correct page metadata and handles page navigation

Test Suites: 1 passed, 1 total
Tests:       5 passed, 5 total
Snapshots:   0 total
```
