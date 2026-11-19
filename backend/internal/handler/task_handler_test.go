package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"task-app/internal/model"
	"task-app/internal/repository"
	"task-app/internal/service"
)

type mockService struct {
	listTasksFn  func(ctx context.Context, filter model.TaskFilter) (*model.TaskListResponse, error)
	getTaskFn    func(ctx context.Context, id int) (*model.Task, error)
	createTaskFn func(ctx context.Context, req *model.CreateTaskRequest) (*model.Task, error)
	updateTaskFn func(ctx context.Context, id int, req *model.UpdateTaskRequest) (*model.Task, error)
	deleteTaskFn func(ctx context.Context, id int) error
}

func (m *mockService) ListTasks(ctx context.Context, filter model.TaskFilter) (*model.TaskListResponse, error) {
	if m.listTasksFn != nil {
		return m.listTasksFn(ctx, filter)
	}
	return nil, nil
}

func (m *mockService) GetTask(ctx context.Context, id int) (*model.Task, error) {
	if m.getTaskFn != nil {
		return m.getTaskFn(ctx, id)
	}
	return nil, nil
}

func (m *mockService) CreateTask(ctx context.Context, req *model.CreateTaskRequest) (*model.Task, error) {
	if m.createTaskFn != nil {
		return m.createTaskFn(ctx, req)
	}
	return nil, nil
}

func (m *mockService) UpdateTask(ctx context.Context, id int, req *model.UpdateTaskRequest) (*model.Task, error) {
	if m.updateTaskFn != nil {
		return m.updateTaskFn(ctx, id, req)
	}
	return nil, nil
}

func (m *mockService) DeleteTask(ctx context.Context, id int) error {
	if m.deleteTaskFn != nil {
		return m.deleteTaskFn(ctx, id)
	}
	return nil
}

func setupTestRouter(svc service.TaskService) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	h := NewTaskHandler(svc)

	api := r.Group("/api/tasks")
	{
		api.GET("", h.GetTasks)
		api.POST("", h.CreateTask)
		api.GET("/:id", h.GetTaskByID)
		api.PUT("/:id", h.UpdateTask)
		api.DELETE("/:id", h.DeleteTask)
	}

	return r
}

// Task 4: Duplicate title should return HTTP 409 instead of 500
func TestHandler_CreateTask_DuplicateTitle_Returns409(t *testing.T) {
	mockSvc := &mockService{
		createTaskFn: func(ctx context.Context, req *model.CreateTaskRequest) (*model.Task, error) {
			return nil, repository.ErrDuplicateTitle
		},
	}

	router := setupTestRouter(mockSvc)

	body, _ := json.Marshal(model.CreateTaskRequest{
		Title: "Existing Title",
	})
	req, _ := http.NewRequest(http.MethodPost, "/api/tasks", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusConflict, w.Code)

	var resp model.APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.False(t, resp.Success)
	assert.NotNil(t, resp.Error)
	assert.Equal(t, "CONFLICT", resp.Error.Code)
	assert.Equal(t, "A task with this title already exists", resp.Error.Message)
}

// Task 1 & 4: PUT /api/tasks/:id with duplicate title returns 409
func TestHandler_UpdateTask_DuplicateTitle_Returns409(t *testing.T) {
	mockSvc := &mockService{
		updateTaskFn: func(ctx context.Context, id int, req *model.UpdateTaskRequest) (*model.Task, error) {
			return nil, repository.ErrDuplicateTitle
		},
	}

	router := setupTestRouter(mockSvc)

	title := "Duplicated Title"
	body, _ := json.Marshal(model.UpdateTaskRequest{
		Title: &title,
	})
	req, _ := http.NewRequest(http.MethodPut, "/api/tasks/1", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusConflict, w.Code)
	var resp model.APIResponse
	_ = json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, "CONFLICT", resp.Error.Code)
}

// Task 1: PUT /api/tasks/:id not found returns 404
func TestHandler_UpdateTask_NotFound_Returns404(t *testing.T) {
	mockSvc := &mockService{
		updateTaskFn: func(ctx context.Context, id int, req *model.UpdateTaskRequest) (*model.Task, error) {
			return nil, repository.ErrNotFound
		},
	}

	router := setupTestRouter(mockSvc)

	title := "Non existent"
	body, _ := json.Marshal(model.UpdateTaskRequest{
		Title: &title,
	})
	req, _ := http.NewRequest(http.MethodPut, "/api/tasks/999", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	var resp model.APIResponse
	_ = json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, "NOT_FOUND", resp.Error.Code)
}

// Task 1: Invalid status filter returns 400
func TestHandler_GetTasks_InvalidStatus_Returns400(t *testing.T) {
	mockSvc := &mockService{}
	router := setupTestRouter(mockSvc)

	req, _ := http.NewRequest(http.MethodGet, "/api/tasks?status=unknown_status", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	var resp model.APIResponse
	_ = json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, "BAD_REQUEST", resp.Error.Code)
}
