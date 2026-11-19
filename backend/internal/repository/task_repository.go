package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/go-sql-driver/mysql"
	"task-app/internal/model"
)

var (
	ErrNotFound       = errors.New("task not found")
	ErrDuplicateTitle = errors.New("task title already exists")
)

type TaskRepository interface {
	Create(ctx context.Context, task *model.Task) (*model.Task, error)
	GetByID(ctx context.Context, id int) (*model.Task, error)
	List(ctx context.Context, filter model.TaskFilter) ([]model.Task, int, error)
	Update(ctx context.Context, id int, task *model.Task) (*model.Task, error)
	SoftDelete(ctx context.Context, id int) error
}

type mysqlTaskRepository struct {
	db *sql.DB
}

func NewMySQLTaskRepository(db *sql.DB) TaskRepository {
	return &mysqlTaskRepository{db: db}
}

func isDuplicateKeyError(err error) bool {
	var mysqlErr *mysql.MySQLError
	if errors.As(err, &mysqlErr) {
		return mysqlErr.Number == 1062
	}
	// Also check error message for sqlite/generic duplicate entries if needed
	if err != nil && strings.Contains(strings.ToLower(err.Error()), "duplicate") {
		return true
	}
	return false
}

func (r *mysqlTaskRepository) Create(ctx context.Context, task *model.Task) (*model.Task, error) {
	query := `INSERT INTO tasks (title, description, status, assignee) VALUES (?, ?, ?, ?)`
	res, err := r.db.ExecContext(ctx, query, task.Title, task.Description, task.Status, task.Assignee)
	if err != nil {
		if isDuplicateKeyError(err) {
			return nil, ErrDuplicateTitle
		}
		return nil, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}

	return r.GetByID(ctx, int(id))
}

func (r *mysqlTaskRepository) GetByID(ctx context.Context, id int) (*model.Task, error) {
	query := `SELECT id, title, description, status, assignee, created_at, updated_at, deleted_at 
	          FROM tasks 
	          WHERE id = ? AND deleted_at IS NULL`

	row := r.db.QueryRowContext(ctx, query, id)

	var t model.Task
	var desc, assignee sql.NullString
	var deletedAt sql.NullTime

	err := row.Scan(&t.ID, &t.Title, &desc, &t.Status, &assignee, &t.CreatedAt, &t.UpdatedAt, &deletedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	if desc.Valid {
		t.Description = &desc.String
	}
	if assignee.Valid {
		t.Assignee = &assignee.String
	}
	if deletedAt.Valid {
		t.DeletedAt = &deletedAt.Time
	}

	return &t, nil
}

func sanitizeSort(sort string) string {
	if strings.TrimSpace(sort) == "" {
		return "created_at DESC"
	}

	parts := strings.Fields(sort)
	col := strings.ToLower(parts[0])
	dir := "ASC"
	if len(parts) > 1 && strings.ToUpper(parts[1]) == "DESC" {
		dir = "DESC"
	}

	allowedCols := map[string]string{
		"id":         "id",
		"title":      "title",
		"status":     "status",
		"assignee":   "assignee",
		"created_at": "created_at",
		"updated_at": "updated_at",
	}

	if validCol, ok := allowedCols[col]; ok {
		return fmt.Sprintf("%s %s", validCol, dir)
	}

	return "created_at DESC"
}

func (r *mysqlTaskRepository) List(ctx context.Context, filter model.TaskFilter) ([]model.Task, int, error) {
	// Base conditions: always hide soft-deleted tasks
	whereClauses := []string{"deleted_at IS NULL"}
	args := []any{}

	if strings.TrimSpace(filter.Status) != "" {
		whereClauses = append(whereClauses, "status = ?")
		args = append(args, strings.TrimSpace(filter.Status))
	}

	if strings.TrimSpace(filter.Assignee) != "" {
		whereClauses = append(whereClauses, "assignee = ?")
		args = append(args, strings.TrimSpace(filter.Assignee))
	}

	if strings.TrimSpace(filter.Keyword) != "" {
		kw := "%" + strings.TrimSpace(filter.Keyword) + "%"
		whereClauses = append(whereClauses, "(title LIKE ? OR description LIKE ?)")
		args = append(args, kw, kw)
	}

	whereSQL := " WHERE " + strings.Join(whereClauses, " AND ")

	// Count total matching items
	countQuery := "SELECT COUNT(*) FROM tasks" + whereSQL
	var total int
	err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// Pagination
	page := filter.Page
	if page < 1 {
		page = 1
	}
	limit := filter.Limit
	if limit < 1 {
		limit = 10
	} else if limit > 100 {
		limit = 100
	}
	offset := (page - 1) * limit

	orderBy := sanitizeSort(filter.Sort)

	dataQuery := fmt.Sprintf(
		"SELECT id, title, description, status, assignee, created_at, updated_at, deleted_at FROM tasks%s ORDER BY %s LIMIT ? OFFSET ?",
		whereSQL,
		orderBy,
	)

	dataArgs := append(args, limit, offset)
	rows, err := r.db.QueryContext(ctx, dataQuery, dataArgs...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	tasks := make([]model.Task, 0)
	for rows.Next() {
		var t model.Task
		var desc, assignee sql.NullString
		var deletedAt sql.NullTime

		if err := rows.Scan(&t.ID, &t.Title, &desc, &t.Status, &assignee, &t.CreatedAt, &t.UpdatedAt, &deletedAt); err != nil {
			return nil, 0, err
		}

		if desc.Valid {
			t.Description = &desc.String
		}
		if assignee.Valid {
			t.Assignee = &assignee.String
		}
		if deletedAt.Valid {
			t.DeletedAt = &deletedAt.Time
		}

		tasks = append(tasks, t)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	return tasks, total, nil
}

func (r *mysqlTaskRepository) Update(ctx context.Context, id int, task *model.Task) (*model.Task, error) {
	// Verify task exists and is not soft-deleted
	existing, err := r.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	query := `UPDATE tasks SET title = ?, description = ?, status = ?, assignee = ?, updated_at = ? 
	          WHERE id = ? AND deleted_at IS NULL`
	now := time.Now()

	res, err := r.db.ExecContext(ctx, query, task.Title, task.Description, task.Status, task.Assignee, now, id)
	if err != nil {
		if isDuplicateKeyError(err) {
			return nil, ErrDuplicateTitle
		}
		return nil, err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return nil, err
	}
	if rowsAffected == 0 {
		// If rows affected is 0, verify if it was no changes or deleted concurrently
		current, err := r.GetByID(ctx, id)
		if err != nil {
			return nil, err
		}
		return current, nil
	}

	existing.Title = task.Title
	existing.Description = task.Description
	existing.Status = task.Status
	existing.Assignee = task.Assignee
	existing.UpdatedAt = now

	return existing, nil
}

func (r *mysqlTaskRepository) SoftDelete(ctx context.Context, id int) error {
	query := `UPDATE tasks SET deleted_at = CURRENT_TIMESTAMP WHERE id = ? AND deleted_at IS NULL`
	res, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return ErrNotFound
	}

	return nil
}
