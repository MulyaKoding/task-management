package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"task-app/internal/model"
	"task-app/internal/repository"
	"task-app/internal/service"
)

func SendSuccess(c *gin.Context, statusCode int, data any, meta any) {
	c.JSON(statusCode, gin.H{
		"success": true,
		"data":    data,
		"meta":    meta,
	})
}

func SendMessage(c *gin.Context, statusCode int, message string) {
	c.JSON(statusCode, gin.H{
		"success": true,
		"message": message,
	})
}

func SendError(c *gin.Context, statusCode int, code string, message string, details any) {
	c.JSON(statusCode, model.APIResponse{
		Success: false,
		Error: &model.APIError{
			Code:    code,
			Message: message,
			Details: details,
		},
	})
}

func HandleServiceError(c *gin.Context, err error) {
	if errors.Is(err, repository.ErrNotFound) {
		SendError(c, http.StatusNotFound, "NOT_FOUND", "Task not found", nil)
		return
	}

	if errors.Is(err, repository.ErrDuplicateTitle) {
		SendError(c, http.StatusConflict, "CONFLICT", "A task with this title already exists", nil)
		return
	}

	if errors.Is(err, service.ErrInvalidStatus) || errors.Is(err, service.ErrEmptyTitle) {
		SendError(c, http.StatusBadRequest, "BAD_REQUEST", err.Error(), nil)
		return
	}

	SendError(c, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "Internal server error", err.Error())
}
