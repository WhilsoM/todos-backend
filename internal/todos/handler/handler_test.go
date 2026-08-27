package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"todos/internal/items"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type fakeService struct {
	todos []items.TodoItem
	err   error
}

func (f *fakeService) GetTodos(ctx context.Context) ([]items.TodoItem, error) {
	return f.todos, f.err
}

func (f *fakeService) CreateTodo(
	ctx context.Context,
	todo items.TodoItem,
) (items.TodoItem, error) {
	return items.TodoItem{ID: 1, Title: todo.Title}, nil
}

func (f *fakeService) DeleteTodoByID(
	ctx context.Context,
	id int,
) error {
	return nil
}

func (f *fakeService) UpdateTodoByID(
	ctx context.Context,
	todo items.TodoItem,
) (items.TodoItem, error) {
	return items.TodoItem{}, nil
}

func TestGetTodos(t *testing.T) {
	service := &fakeService{
		todos: []items.TodoItem{
			{
				ID:    1,
				Title: "learn go",
			},
		},
	}

	handler := NewHandler(service)

	req := httptest.NewRequest(http.MethodGet, "/api/todos", nil)
	rec := httptest.NewRecorder()

	handler.GetTodos(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("got status: %d, want: %d", rec.Code, http.StatusOK)
	}

	todos := make([]items.TodoItem, 0)

	require.NoError(t, json.NewDecoder(rec.Body).Decode(&todos))

	require.Len(t, todos, 1)
	assert.Equal(t, 1, todos[0].ID)
	assert.Equal(t, "learn go", todos[0].Title)
}
