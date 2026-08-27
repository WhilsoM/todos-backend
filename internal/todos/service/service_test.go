package service

import (
	"context"
	"log"
	"testing"
	customerrors "todos/internal/custom-errors"
	"todos/internal/items"

	"github.com/stretchr/testify/require"
)

type fakeRepository struct {
	todos []items.TodoItem
	err   error

	createResult items.TodoItem
	createErr    error
}

func (f *fakeRepository) GetTodos(ctx context.Context) ([]items.TodoItem, error) {
	return f.todos, f.err
}

func (f *fakeRepository) CreateTodo(
	ctx context.Context,
	todo items.TodoItem,
) (items.TodoItem, error) {
	return f.createResult, f.createErr
}

func (f *fakeRepository) DeleteTodoByID(
	ctx context.Context,
	id int,
) error {
	return nil
}

func (f *fakeRepository) UpdateTodoByID(
	ctx context.Context,
	todo items.TodoItem,
) (items.TodoItem, error) {
	return items.TodoItem{}, nil
}

func TestGetTodos(t *testing.T) {
	expectedTodos := []items.TodoItem{
		{
			ID:    1,
			Title: "learn go",
		},
	}

	repo := &fakeRepository{
		todos: expectedTodos,
	}

	service := NewService(repo, nil)

	ctx := context.Background()

	todos, err := service.GetTodos(ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(todos) != 1 {
		t.Fatalf("got length: %d want: %d", len(todos), 1)
	}
}

func TestCreateTodo(t *testing.T) {
	tests := []struct {
		name         string
		title        string
		createResult items.TodoItem
		wantErr      error
		wantTitle    string
	}{
		{
			name:  "valid title",
			title: "learn go",

			createResult: items.TodoItem{
				ID:    1,
				Title: "learn go",
			},

			wantErr:   nil,
			wantTitle: "learn go",
		},
		{
			name:      "empty title",
			title:     "",
			wantErr:   customerrors.ErrTitleEmpty,
			wantTitle: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeRepo := &fakeRepository{
				createResult: tt.createResult,
			}
			ctx := context.Background()

			todo := items.TodoItem{
				Title: tt.title,
			}

			service := NewService(fakeRepo, nil)

			newTodo, err := service.CreateTodo(ctx, todo)
			if err != tt.wantErr {
				require.ErrorIs(t, err, tt.wantErr)
				t.Fatalf("unexpected error got: %v, want: %v", err, tt.wantErr)
			}

			log.Println(newTodo.Title, tt.title)
			if newTodo.Title != tt.title {
				t.Fatalf("title cannot be not equal got: %s, want: %s", newTodo.Title, tt.title)
			}
		})
	}
}
