package redis

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"time"
	"todos/internal/items"

	"github.com/redis/go-redis/v9"
)

type TodoCache struct {
	client *redis.Client
	ttl    time.Duration
}

func NewTodoCache(client *redis.Client) *TodoCache {
	return &TodoCache{
		client: client,
		ttl:    30 * time.Second,
	}
}

func (t *TodoCache) GetTodos(ctx context.Context) ([]items.TodoItem, bool, error) {
	slog.Info("GET TODOS CACHE STARTED")
	val, err := t.client.Get(ctx, "todos:all").Result()
	if errors.Is(err, redis.Nil) {
		slog.Info("GET TODOS CACHE NOT FOUND")

		return nil, false, nil
	}

	if err != nil {
		slog.Info("GET TODOS CACHE error")
		return nil, false, err
	}

	var results []items.TodoItem
	if err := json.Unmarshal([]byte(val), &results); err != nil {
		return nil, false, err
	}
	slog.Info("GET TODOS CACHE HIT")

	return results, true, nil
}

func (t *TodoCache) SetTodos(ctx context.Context, todos []items.TodoItem) error {
	val, err := json.Marshal(todos)
	if err != nil {
		return err
	}

	if err := t.client.Set(ctx, "todos:all", val, t.ttl).Err(); err != nil {
		return err
	}

	return nil
}

func (t *TodoCache) DeleteTodos(ctx context.Context) error {
	return t.client.Del(ctx, "todos:all").Err()
}
