package cache

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
	"task-app/internal/model"
)

const (
	TaskCachePrefix = "tasks:"
	TaskCacheTTL    = 60 * time.Second
)

type TaskCache interface {
	Get(ctx context.Context, key string) ([]byte, error)
	Set(ctx context.Context, key string, data []byte, ttl time.Duration) error
	InvalidateAll(ctx context.Context) error
	BuildKey(filter model.TaskFilter) string
}

type redisTaskCache struct {
	client *redis.Client
}

func NewRedisTaskCache(client *redis.Client) TaskCache {
	return &redisTaskCache{client: client}
}

func (c *redisTaskCache) BuildKey(filter model.TaskFilter) string {
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
	sort := strings.TrimSpace(filter.Sort)
	if sort == "" {
		sort = "created_at DESC"
	}

	return fmt.Sprintf(
		"%sstatus=%s:keyword=%s:assignee=%s:page=%d:limit=%d:sort=%s",
		TaskCachePrefix,
		strings.TrimSpace(filter.Status),
		strings.TrimSpace(filter.Keyword),
		strings.TrimSpace(filter.Assignee),
		page,
		limit,
		sort,
	)
}

func (c *redisTaskCache) Get(ctx context.Context, key string) ([]byte, error) {
	if c.client == nil {
		return nil, nil
	}

	val, err := c.client.Get(ctx, key).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		return nil, err
	}
	return val, nil
}

func (c *redisTaskCache) Set(ctx context.Context, key string, data []byte, ttl time.Duration) error {
	if c.client == nil {
		return nil
	}
	return c.client.Set(ctx, key, data, ttl).Err()
}

func (c *redisTaskCache) InvalidateAll(ctx context.Context) error {
	if c.client == nil {
		return nil
	}

	pattern := TaskCachePrefix + "*"
	var cursor uint64
	var totalDeleted int64

	for {
		var keys []string
		var err error
		keys, cursor, err = c.client.Scan(ctx, cursor, pattern, 50).Result()
		if err != nil {
			log.Printf("[Redis] Scan error during invalidation: %v", err)
			return err
		}

		if len(keys) > 0 {
			deleted, err := c.client.Del(ctx, keys...).Result()
			if err != nil {
				log.Printf("[Redis] Del error during invalidation: %v", err)
				return err
			}
			totalDeleted += deleted
		}

		if cursor == 0 {
			break
		}
	}

	log.Printf("[Redis] Successfully invalidated %d cache keys matching '%s'", totalDeleted, pattern)
	return nil
}
