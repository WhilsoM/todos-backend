package repository

import (
	"context"
	"errors"
	customerrors "todos/internal/custom-errors"
	"todos/internal/items"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{
		db,
	}
}

func (r *Repository) GetTodos(ctx context.Context) ([]items.TodoItem, error) {
	rows, err := r.db.Query(ctx, `SELECT id, title, success FROM todos`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	todos := make([]items.TodoItem, 0)

	for rows.Next() {
		var todo items.TodoItem

		err := rows.Scan(&todo.ID, &todo.Title, &todo.Success)
		if err != nil {
			return nil, err
		}
		todos = append(todos, todo)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return todos, nil
}

func (r *Repository) CreateTodo(ctx context.Context, todo items.TodoItem) (items.TodoItem, error) {
	query := `INSERT INTO todos (title, success) VALUES ($1, false) RETURNING id`
	id := 0

	row := r.db.QueryRow(ctx, query, todo.Title)

	if err := row.Scan(&id); err != nil {
		return items.TodoItem{}, err
	}

	return items.TodoItem{
		ID:    id,
		Title: todo.Title,
	}, nil

}

func (r *Repository) DeleteTodoByID(ctx context.Context, id int) error {
	query := `DELETE FROM todos WHERE id = $1`

	command, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	if command.RowsAffected() == 0 {
		return customerrors.ErrNotFound
	}

	return nil
}

func (r *Repository) UpdateTodoByID(ctx context.Context, newTodo items.TodoItem) (items.TodoItem, error) {
	query := `UPDATE todos SET title = $1, success = $2 WHERE id = $3 RETURNING title, success, id`

	var returnedTodo items.TodoItem

	row := r.db.QueryRow(ctx, query, newTodo.Title, newTodo.Success, newTodo.ID)
	if err := row.Scan(&returnedTodo.Title, &returnedTodo.Success, &returnedTodo.ID); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return items.TodoItem{}, customerrors.ErrNotFound
		}
		return items.TodoItem{}, err
	}

	return returnedTodo, nil
}
