package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"task-app/internal/cache"
	"task-app/internal/model"
	"task-app/internal/repository"
)

var (
	ErrInvalidStatus = errors.New("status must be one of: todo, in_progress, done")
	ErrEmptyTitle    = errors.New("title cannot be empty")
)

type TaskService interface {
	ListTasks(ctx context.Context, filter model.TaskFilter) (*model.TaskListResponse, error)
	GetTask(ctx context.Context, id int) (*model.Task, error)
	CreateTask(ctx context.Context, req *model.CreateTaskRequest) (*model.Task, error)
	UpdateTask(ctx context.Context, id int, req *model.UpdateTaskRequest) (*model.Task, error)
	DeleteTask(ctx context.Context, id int) error
}

type taskService struct {
	repo  repository.TaskRepository
	cache cache.TaskCache
}

func NewTaskService(repo repository.TaskRepository, cache cache.TaskCache) TaskService {
	return &taskService{
		repo:  repo,
		cache: cache,
	}
}

func (s *taskService) ListTasks(ctx context.Context, filter model.TaskFilter) (*model.TaskListResponse, error) {
	// Normalize filter parameters
	if filter.Page < 1 {
		filter.Page = 1
	}
	if filter.Limit < 1 {
		filter.Limit = 10
	} else if filter.Limit > 100 {
		filter.Limit = 100
	}

	cacheKey := s.cache.BuildKey(filter)

	// Try reading from cache
	cached, err := s.cache.Get(ctx, cacheKey)
	if err == nil && cached != nil {
		var resp model.TaskListResponse
		if err := json.Unmarshal(cached, &resp); err == nil {
			log.Printf("[Cache HIT] Key: %s", cacheKey)
			return &resp, nil
		}
	}
	log.Printf("[Cache MISS] Key: %s", cacheKey)

	// Fetch from database
	tasks, total, err := s.repo.List(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("repository list failed: %w", err)
	}

	totalPages := 0
	if total > 0 {
		totalPages = (total + filter.Limit - 1) / filter.Limit
	}

	response := &model.TaskListResponse{
		Data: tasks,
		Meta: model.PaginationMeta{
			Page:       filter.Page,
			Limit:      filter.Limit,
			Total:      total,
			TotalPages: totalPages,
		},
	}

	// Cache the result for 60 seconds
	if data, err := json.Marshal(response); err == nil {
		if err := s.cache.Set(ctx, cacheKey, data, cache.TaskCacheTTL); err != nil {
			log.Printf("[Cache SET ERROR] Failed setting key %s: %v", cacheKey, err)
		}
	}

	return response, nil
}

func (s *taskService) GetTask(ctx context.Context, id int) (*model.Task, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *taskService) CreateTask(ctx context.Context, req *model.CreateTaskRequest) (*model.Task, error) {
	if req.Title == "" {
		return nil, ErrEmptyTitle
	}

	status := model.StatusTodo
	if req.Status != nil {
		if !req.Status.IsValid() {
			return nil, ErrInvalidStatus
		}
		status = *req.Status
	}

	task := &model.Task{
		Title:       req.Title,
		Description: req.Description,
		Status:      status,
		Assignee:    req.Assignee,
	}

	created, err := s.repo.Create(ctx, task)
	if err != nil {
		return nil, err
	}

	// Invalidate cache after create
	if err := s.cache.InvalidateAll(ctx); err != nil {
		log.Printf("[Cache Invalidate Warning] Failed to invalidate after create: %v", err)
	}

	return created, nil
}

func (s *taskService) UpdateTask(ctx context.Context, id int, req *model.UpdateTaskRequest) (*model.Task, error) {
	existing, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if req.Title != nil {
		if *req.Title == "" {
			return nil, ErrEmptyTitle
		}
		existing.Title = *req.Title
	}

	if req.Description != nil {
		existing.Description = req.Description
	}

	if req.Status != nil {
		if !req.Status.IsValid() {
			return nil, ErrInvalidStatus
		}
		existing.Status = *req.Status
	}

	if req.Assignee != nil {
		existing.Assignee = req.Assignee
	}

	updated, err := s.repo.Update(ctx, id, existing)
	if err != nil {
		return nil, err
	}

	// Invalidate cache after update
	if err := s.cache.InvalidateAll(ctx); err != nil {
		log.Printf("[Cache Invalidate Warning] Failed to invalidate after update: %v", err)
	}

	return updated, nil
}

func (s *taskService) DeleteTask(ctx context.Context, id int) error {
	if err := s.repo.SoftDelete(ctx, id); err != nil {
		return err
	}

	// Invalidate cache after soft-delete
	if err := s.cache.InvalidateAll(ctx); err != nil {
		log.Printf("[Cache Invalidate Warning] Failed to invalidate after delete: %v", err)
	}

	return nil
}
