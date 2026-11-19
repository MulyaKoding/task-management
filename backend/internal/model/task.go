package model

import "time"

type TaskStatus string

const (
	StatusTodo       TaskStatus = "todo"
	StatusInProgress TaskStatus = "in_progress"
	StatusDone       TaskStatus = "done"
)

func (s TaskStatus) IsValid() bool {
	switch s {
	case StatusTodo, StatusInProgress, StatusDone:
		return true
	default:
		return false
	}
}

type Task struct {
	ID          int        `json:"id"`
	Title       string     `json:"title"`
	Description *string    `json:"description"`
	Status      TaskStatus `json:"status"`
	Assignee    *string    `json:"assignee"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty"`
}

type CreateTaskRequest struct {
	Title       string      `json:"title" binding:"required"`
	Description *string     `json:"description"`
	Status      *TaskStatus `json:"status"`
	Assignee    *string     `json:"assignee"`
}

type UpdateTaskRequest struct {
	Title       *string     `json:"title"`
	Description *string     `json:"description"`
	Status      *TaskStatus `json:"status"`
	Assignee    *string     `json:"assignee"`
}

type TaskFilter struct {
	Status   string `form:"status"`
	Keyword  string `form:"keyword"`
	Assignee string `form:"assignee"`
	Page     int    `form:"page"`
	Limit    int    `form:"limit"`
	Sort     string `form:"sort"`
}

type PaginationMeta struct {
	Page       int `json:"page"`
	Limit      int `json:"limit"`
	Total      int `json:"total"`
	TotalPages int `json:"total_pages"`
}

type TaskListResponse struct {
	Data []Task         `json:"data"`
	Meta PaginationMeta `json:"meta"`
}

type APIError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details any    `json:"details,omitempty"`
}

type APIResponse struct {
	Success bool      `json:"success"`
	Data    any       `json:"data,omitempty"`
	Meta    any       `json:"meta,omitempty"`
	Error   *APIError `json:"error,omitempty"`
}
