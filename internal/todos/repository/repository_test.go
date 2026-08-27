package repository

import (
	"context"
	"testing"
	"todos/internal/items"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateTodo(t *testing.T) {
	tests := []struct {
		name      string
		title     string
		wantTitle string
	}{
		{
			name:      "learn go",
			title:     "learn go",
			wantTitle: "learn go",
		},
		{
			name:      "learn english",
			title:     "learn english",
			wantTitle: "learn english",
		},
		{
			name:      "learn math",
			title:     "learn math",
			wantTitle: "learn math",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := NewRepository(nil)
			ctx := context.Background()

			todo := items.TodoItem{
				Title: tt.title,
			}

			newTodo, err := repo.CreateTodo(ctx, todo)
			require.NoError(t, err)

			assert.Equal(t, 1, newTodo.ID)
			assert.Equal(t, tt.wantTitle, newTodo.Title)
		})
	}

}

func TestGetTodos_EmptyArray(t *testing.T) {
	repo := NewRepository(nil)
	ctx := context.Background()

	todos, err := repo.GetTodos(ctx)
	require.NoError(t, err)

	assert.Len(t, todos, 0)
}

func TestGetTodos_SomeItemsArray(t *testing.T) {
	repo := NewRepository(nil)
	ctx := context.Background()

	todo1 := items.TodoItem{
		Title: "learn go",
	}
	todo2 := items.TodoItem{
		Title: "learn english",
	}

	_, err := repo.CreateTodo(ctx, todo1)
	require.NoError(t, err)

	_, err = repo.CreateTodo(ctx, todo2)
	require.NoError(t, err)

	todos, err := repo.GetTodos(ctx)
	require.NoError(t, err)

	assert.Len(t, todos, 2)
}

func TestDeleteTodoByID(t *testing.T) {
	repo := NewRepository(nil)
	ctx := context.Background()

	todo := items.TodoItem{
		Title: "learn go",
	}

	newTodo, err := repo.CreateTodo(ctx, todo)
	require.NoError(t, err)

	err = repo.DeleteTodoByID(ctx, newTodo.ID)
	require.NoError(t, err)

	todos, err := repo.GetTodos(ctx)
	require.NoError(t, err)

	assert.Len(t, todos, 0)
}

func TestUpdateTodoByID(t *testing.T) {
	repo := NewRepository(nil)
	ctx := context.Background()

	todo := items.TodoItem{
		Title: "learn golang",
	}

	newTodo, err := repo.CreateTodo(ctx, todo)
	require.NoError(t, err)

	newTodo2 := items.TodoItem{
		ID:    newTodo.ID,
		Title: "learn math",
	}

	updatedTodo, err := repo.UpdateTodoByID(ctx, newTodo2)
	require.NoError(t, err)

	assert.Equal(t, "learn math", updatedTodo.Title)
}
