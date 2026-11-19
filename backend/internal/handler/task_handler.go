package handler

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"task-app/internal/model"
	"task-app/internal/service"
)

type TaskHandler struct {
	svc service.TaskService
}

func NewTaskHandler(svc service.TaskService) *TaskHandler {
	return &TaskHandler{svc: svc}
}

// GetTasks handles GET /api/tasks with filtering, pagination, and sorting
func (h *TaskHandler) GetTasks(c *gin.Context) {
	var filter model.TaskFilter

	if err := c.ShouldBindQuery(&filter); err != nil {
		SendError(c, http.StatusBadRequest, "BAD_REQUEST", "Invalid query parameters", err.Error())
		return
	}

	if filter.Status != "" {
		st := model.TaskStatus(filter.Status)
		if !st.IsValid() {
			SendError(c, http.StatusBadRequest, "BAD_REQUEST", "Invalid status. Allowed values: todo, in_progress, done", nil)
			return
		}
	}

	if filter.Page < 1 {
		filter.Page = 1
	}
	if filter.Limit < 1 {
		filter.Limit = 10
	}

	result, err := h.svc.ListTasks(c.Request.Context(), filter)
	if err != nil {
		HandleServiceError(c, err)
		return
	}

	SendSuccess(c, http.StatusOK, result.Data, result.Meta)
}

// GetTaskByID handles GET /api/tasks/:id
func (h *TaskHandler) GetTaskByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil || id <= 0 {
		SendError(c, http.StatusBadRequest, "BAD_REQUEST", "Invalid task ID, must be a positive integer", nil)
		return
	}

	task, err := h.svc.GetTask(c.Request.Context(), id)
	if err != nil {
		HandleServiceError(c, err)
		return
	}

	SendSuccess(c, http.StatusOK, task, nil)
}

// CreateTask handles POST /api/tasks
func (h *TaskHandler) CreateTask(c *gin.Context) {
	var req model.CreateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		SendError(c, http.StatusBadRequest, "BAD_REQUEST", "Invalid request body", err.Error())
		return
	}

	if strings.TrimSpace(req.Title) == "" {
		SendError(c, http.StatusBadRequest, "BAD_REQUEST", "Title is required", nil)
		return
	}

	task, err := h.svc.CreateTask(c.Request.Context(), &req)
	if err != nil {
		HandleServiceError(c, err)
		return
	}

	SendSuccess(c, http.StatusCreated, task, nil)
}

// UpdateTask handles PUT /api/tasks/:id
func (h *TaskHandler) UpdateTask(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil || id <= 0 {
		SendError(c, http.StatusBadRequest, "BAD_REQUEST", "Invalid task ID, must be a positive integer", nil)
		return
	}

	var req model.UpdateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		SendError(c, http.StatusBadRequest, "BAD_REQUEST", "Invalid request body", err.Error())
		return
	}

	task, err := h.svc.UpdateTask(c.Request.Context(), id, &req)
	if err != nil {
		HandleServiceError(c, err)
		return
	}

	SendSuccess(c, http.StatusOK, task, nil)
}

// DeleteTask handles DELETE /api/tasks/:id (soft delete)
func (h *TaskHandler) DeleteTask(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil || id <= 0 {
		SendError(c, http.StatusBadRequest, "BAD_REQUEST", "Invalid task ID, must be a positive integer", nil)
		return
	}

	if err := h.svc.DeleteTask(c.Request.Context(), id); err != nil {
		HandleServiceError(c, err)
		return
	}

	SendMessage(c, http.StatusOK, "Task successfully soft-deleted")
}
