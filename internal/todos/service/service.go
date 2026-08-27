package service

import (
	"context"
	"log/slog"
	"strings"
	customerrors "todos/internal/custom-errors"
	"todos/internal/items"
)

type TodosRedis interface {
	GetTodos(ctx context.Context) ([]items.TodoItem, error)
}

type TodoCache interface {
	GetTodos(ctx context.Context) ([]items.TodoItem, bool, error)
	SetTodos(ctx context.Context, todos []items.TodoItem) error
	DeleteTodos(ctx context.Context) error
}

type RepositoryInterface interface {
	GetTodos(ctx context.Context) ([]items.TodoItem, error)
	CreateTodo(ctx context.Context, todo items.TodoItem) (items.TodoItem, error)
	DeleteTodoByID(ctx context.Context, id int) error
	UpdateTodoByID(ctx context.Context, newTodo items.TodoItem) (items.TodoItem, error)
}

type Service struct {
	repo        RepositoryInterface
	redisClient TodoCache
}

func NewService(repo RepositoryInterface, redisClient TodoCache) *Service {
	return &Service{
		repo,
		redisClient,
	}
}

func (s *Service) GetTodos(ctx context.Context) ([]items.TodoItem, error) {
	todosCache, found, err := s.redisClient.GetTodos(ctx)
	if found {
		return todosCache, nil
	}

	todos, err := s.repo.GetTodos(ctx)
	if err != nil {
		return todos, err
	}

	if err := s.redisClient.SetTodos(ctx, todos); err != nil {
		slog.Info("failed to save todos to cache", "err", err, "todos", todos)
	}

	return todos, nil
}

func (s *Service) CreateTodo(ctx context.Context, todo items.TodoItem) (items.TodoItem, error) {
	if strings.TrimSpace(todo.Title) == "" {
		return items.TodoItem{}, customerrors.ErrTitleEmpty
	}

	todo, err := s.repo.CreateTodo(ctx, todo)
	if err != nil {
		return todo, err
	}

	return todo, nil
}

func (s *Service) DeleteTodoByID(ctx context.Context, id int) error {
	if err := s.repo.DeleteTodoByID(ctx, id); err != nil {
		return err
	}

	return nil
}

func (s *Service) UpdateTodoByID(ctx context.Context, newTodo items.TodoItem) (items.TodoItem, error) {
	todo, err := s.repo.UpdateTodoByID(ctx, newTodo)
	if err != nil {
		return todo, err
	}

	return todo, nil
}
