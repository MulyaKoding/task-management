package repository

import (
	"context"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
	"task-app/internal/model"
)

func TestMySQLTaskRepository_Update_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewMySQLTaskRepository(db)
	now := time.Now()

	title := "Updated Task"
	desc := "Updated Description"
	status := model.StatusInProgress
	assignee := "alice"

	// Mock GetByID inside Update
	mock.ExpectQuery(`SELECT id, title, description, status, assignee, created_at, updated_at, deleted_at FROM tasks WHERE id = \? AND deleted_at IS NULL`).
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "title", "description", "status", "assignee", "created_at", "updated_at", "deleted_at"}).
			AddRow(1, "Old Title", "Old Desc", "todo", "alice", now, now, nil))

	// Mock UPDATE query
	mock.ExpectExec(`UPDATE tasks SET title = \?, description = \?, status = \?, assignee = \?, updated_at = \? WHERE id = \? AND deleted_at IS NULL`).
		WithArgs(title, &desc, status, &assignee, sqlmock.AnyArg(), 1).
		WillReturnResult(sqlmock.NewResult(0, 1))

	updated, err := repo.Update(context.Background(), 1, &model.Task{
		Title:       title,
		Description: &desc,
		Status:      status,
		Assignee:    &assignee,
	})

	assert.NoError(t, err)
	assert.NotNil(t, updated)
	assert.Equal(t, title, updated.Title)
	assert.Equal(t, &desc, updated.Description)
	assert.Equal(t, status, updated.Status)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestMySQLTaskRepository_Update_DuplicateTitle_ReturnsConflict(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewMySQLTaskRepository(db)
	now := time.Now()

	title := "Already Existing Title"
	desc := "Desc"

	// Mock GetByID
	mock.ExpectQuery(`SELECT id, title, description, status, assignee, created_at, updated_at, deleted_at FROM tasks WHERE id = \? AND deleted_at IS NULL`).
		WithArgs(2).
		WillReturnRows(sqlmock.NewRows([]string{"id", "title", "description", "status", "assignee", "created_at", "updated_at", "deleted_at"}).
			AddRow(2, "Old Title", "Desc", "todo", nil, now, now, nil))

	// Mock UPDATE throwing MySQL 1062 duplicate entry error
	mysqlErr := &mysql.MySQLError{
		Number:  1062,
		Message: "Duplicate entry 'Already Existing Title' for key 'uq_title'",
	}
	mock.ExpectExec(`UPDATE tasks SET title = \?, description = \?, status = \?, assignee = \?, updated_at = \? WHERE id = \? AND deleted_at IS NULL`).
		WithArgs(title, &desc, model.StatusTodo, (*string)(nil), sqlmock.AnyArg(), 2).
		WillReturnError(mysqlErr)

	updated, err := repo.Update(context.Background(), 2, &model.Task{
		Title:       title,
		Description: &desc,
		Status:      model.StatusTodo,
	})

	assert.Nil(t, updated)
	assert.ErrorIs(t, err, ErrDuplicateTitle)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestMySQLTaskRepository_List_SearchAndHideSoftDeleted(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewMySQLTaskRepository(db)
	now := time.Now()

	filter := model.TaskFilter{
		Status:   "todo",
		Keyword:  "frontend",
		Assignee: "john",
		Page:     1,
		Limit:    10,
	}

	kw := "%frontend%"

	// Mock COUNT query ensuring deleted_at IS NULL
	mock.ExpectQuery(`SELECT COUNT\(\*\) FROM tasks WHERE deleted_at IS NULL AND status = \? AND assignee = \? AND \(title LIKE \? OR description LIKE \?\)`).
		WithArgs("todo", "john", kw, kw).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

	// Mock DATA query ensuring deleted_at IS NULL
	mock.ExpectQuery(`SELECT id, title, description, status, assignee, created_at, updated_at, deleted_at FROM tasks WHERE deleted_at IS NULL AND status = \? AND assignee = \? AND \(title LIKE \? OR description LIKE \?\) ORDER BY created_at DESC LIMIT \? OFFSET \?`).
		WithArgs("todo", "john", kw, kw, 10, 0).
		WillReturnRows(sqlmock.NewRows([]string{"id", "title", "description", "status", "assignee", "created_at", "updated_at", "deleted_at"}).
			AddRow(10, "Fix frontend search", "Search bug in list", "todo", "john", now, now, nil))

	tasks, total, err := repo.List(context.Background(), filter)

	assert.NoError(t, err)
	assert.Equal(t, 1, total)
	assert.Len(t, tasks, 1)
	assert.Equal(t, 10, tasks[0].ID)
	assert.Equal(t, "Fix frontend search", tasks[0].Title)
	assert.Nil(t, tasks[0].DeletedAt) // Confirms soft-deleted tasks are not present
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestMySQLTaskRepository_SoftDelete(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewMySQLTaskRepository(db)

	// Mock SoftDelete query
	mock.ExpectExec(`UPDATE tasks SET deleted_at = CURRENT_TIMESTAMP WHERE id = \? AND deleted_at IS NULL`).
		WithArgs(5).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err = repo.SoftDelete(context.Background(), 5)
	assert.NoError(t, err)

	// SoftDelete on non-existent or already soft-deleted task returns ErrNotFound
	mock.ExpectExec(`UPDATE tasks SET deleted_at = CURRENT_TIMESTAMP WHERE id = \? AND deleted_at IS NULL`).
		WithArgs(999).
		WillReturnResult(sqlmock.NewResult(0, 0))

	err = repo.SoftDelete(context.Background(), 999)
	assert.ErrorIs(t, err, ErrNotFound)
	assert.NoError(t, mock.ExpectationsWereMet())
}
