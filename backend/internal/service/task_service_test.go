package service

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"task-app/internal/model"
)

type mockRepository struct {
	createFn func(ctx context.Context, task *model.Task) (*model.Task, error)
	getFn    func(ctx context.Context, id int) (*model.Task, error)
	listFn   func(ctx context.Context, filter model.TaskFilter) ([]model.Task, int, error)
	updateFn func(ctx context.Context, id int, task *model.Task) (*model.Task, error)
	deleteFn func(ctx context.Context, id int) error
}

func (m *mockRepository) Create(ctx context.Context, task *model.Task) (*model.Task, error) {
	if m.createFn != nil {
		return m.createFn(ctx, task)
	}
	return nil, nil
}

func (m *mockRepository) GetByID(ctx context.Context, id int) (*model.Task, error) {
	if m.getFn != nil {
		return m.getFn(ctx, id)
	}
	return nil, nil
}

func (m *mockRepository) List(ctx context.Context, filter model.TaskFilter) ([]model.Task, int, error) {
	if m.listFn != nil {
		return m.listFn(ctx, filter)
	}
	return nil, 0, nil
}

func (m *mockRepository) Update(ctx context.Context, id int, task *model.Task) (*model.Task, error) {
	if m.updateFn != nil {
		return m.updateFn(ctx, id, task)
	}
	return nil, nil
}

func (m *mockRepository) SoftDelete(ctx context.Context, id int) error {
	if m.deleteFn != nil {
		return m.deleteFn(ctx, id)
	}
	return nil
}

type mockCacheStore struct {
	store          map[string][]byte
	invalidated    bool
	invalidateCall int
	getCall        int
	setCall        int
}

func newMockCacheStore() *mockCacheStore {
	return &mockCacheStore{
		store: make(map[string][]byte),
	}
}

func (m *mockCacheStore) Get(ctx context.Context, key string) ([]byte, error) {
	m.getCall++
	val, ok := m.store[key]
	if !ok {
		return nil, nil
	}
	return val, nil
}

func (m *mockCacheStore) Set(ctx context.Context, key string, data []byte, ttl time.Duration) error {
	m.setCall++
	m.store[key] = data
	return nil
}

func (m *mockCacheStore) InvalidateAll(ctx context.Context) error {
	m.invalidated = true
	m.invalidateCall++
	m.store = make(map[string][]byte)
	return nil
}

func (m *mockCacheStore) BuildKey(filter model.TaskFilter) string {
	return "tasks:mock-key"
}

// Test Task 2 & Task 5: Search with Cache Miss and Cache Hit
func TestTaskService_Search_CacheMissAndHit(t *testing.T) {
	ctx := context.Background()
	c := newMockCacheStore()

	repoCalled := 0
	sampleTask := model.Task{
		ID:        1,
		Title:     "Task Alpha",
		Status:    model.StatusTodo,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	repo := &mockRepository{
		listFn: func(ctx context.Context, filter model.TaskFilter) ([]model.Task, int, error) {
			repoCalled++
			return []model.Task{sampleTask}, 1, nil
		},
	}

	svc := NewTaskService(repo, c)
	filter := model.TaskFilter{Keyword: "Alpha", Page: 1, Limit: 10}

	// 1st Call: Cache MISS -> Must query repository and store into cache
	resp1, err := svc.ListTasks(ctx, filter)
	assert.NoError(t, err)
	assert.NotNil(t, resp1)
	assert.Equal(t, 1, repoCalled)
	assert.Equal(t, 1, c.getCall)
	assert.Equal(t, 1, c.setCall)
	assert.Len(t, resp1.Data, 1)
	assert.Equal(t, "Task Alpha", resp1.Data[0].Title)

	// Verify cached content exists
	cachedData, _ := c.Get(ctx, "tasks:mock-key")
	assert.NotNil(t, cachedData)

	// 2nd Call: Cache HIT -> Must read directly from cache without hitting repo
	resp2, err := svc.ListTasks(ctx, filter)
	assert.NoError(t, err)
	assert.NotNil(t, resp2)
	assert.Equal(t, 1, repoCalled, "Repo should NOT be called again on cache hit")
	assert.Equal(t, 3, c.getCall) // 1st search + manual Get above + 2nd search
	assert.Len(t, resp2.Data, 1)
	assert.Equal(t, "Task Alpha", resp2.Data[0].Title)
}

// Test Task 2 & Task 5: Cache Invalidation on Update
func TestTaskService_Update_InvalidatesCache(t *testing.T) {
	ctx := context.Background()
	c := newMockCacheStore()
	c.store["tasks:some-cached-key"] = []byte(`{"data":[]}`)

	existing := &model.Task{
		ID:     1,
		Title:  "Original Title",
		Status: model.StatusTodo,
	}

	repo := &mockRepository{
		getFn: func(ctx context.Context, id int) (*model.Task, error) {
			return existing, nil
		},
		updateFn: func(ctx context.Context, id int, task *model.Task) (*model.Task, error) {
			return task, nil
		},
	}

	svc := NewTaskService(repo, c)
	newTitle := "Updated Title"
	newStatus := model.StatusInProgress

	updated, err := svc.UpdateTask(ctx, 1, &model.UpdateTaskRequest{
		Title:  &newTitle,
		Status: &newStatus,
	})

	assert.NoError(t, err)
	assert.NotNil(t, updated)
	assert.Equal(t, "Updated Title", updated.Title)
	assert.Equal(t, model.StatusInProgress, updated.Status)
	assert.True(t, c.invalidated, "Cache must be invalidated after task update")
	assert.Equal(t, 1, c.invalidateCall)
}

// Test Task 2 & Task 5: Cache Invalidation on Create & Soft-Delete
func TestTaskService_CreateAndDelete_InvalidatesCache(t *testing.T) {
	ctx := context.Background()
	c := newMockCacheStore()

	repo := &mockRepository{
		createFn: func(ctx context.Context, task *model.Task) (*model.Task, error) {
			task.ID = 100
			return task, nil
		},
		deleteFn: func(ctx context.Context, id int) error {
			return nil
		},
	}

	svc := NewTaskService(repo, c)

	// Create Task
	created, err := svc.CreateTask(ctx, &model.CreateTaskRequest{
		Title: "Brand New Task",
	})
	assert.NoError(t, err)
	assert.NotNil(t, created)
	assert.Equal(t, 1, c.invalidateCall, "Cache must be invalidated after create")

	// Delete Task
	err = svc.DeleteTask(ctx, 100)
	assert.NoError(t, err)
	assert.Equal(t, 2, c.invalidateCall, "Cache must be invalidated after soft delete")
}
