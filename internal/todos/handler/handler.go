package handler

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	customerrors "todos/internal/custom-errors"
	"todos/internal/items"
	"todos/internal/middleware"
	"todos/internal/utils"
)

type ServiceInterface interface {
	GetTodos(ctx context.Context) ([]items.TodoItem, error)
	CreateTodo(ctx context.Context, todo items.TodoItem) (items.TodoItem, error)
	DeleteTodoByID(ctx context.Context, id int) error
	UpdateTodoByID(ctx context.Context, newTodo items.TodoItem) (items.TodoItem, error)
}

type Handler struct {
	service ServiceInterface
}

func NewHandler(service ServiceInterface) *Handler {
	return &Handler{
		service,
	}
}

func (h *Handler) GetTodos(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID, ok := middleware.UserIDFromContext(ctx)
	if !ok {
		slog.Error("failed to get todos user id is not ok", "ok", ok)
		utils.WriteJSONResponseError(w, http.StatusBadRequest, "user id is not correct")
		return
	}
	slog.Info("get todos", "userid", fmt.Sprintf("%d", userID))

	todos, err := h.service.GetTodos(ctx)
	if err != nil {
		slog.Error("failed to get todos", "error", err)
		utils.WriteJSONResponseError(w, http.StatusBadRequest, "error: "+err.Error())
		return
	}

	utils.WriteJSONResponse(w, http.StatusOK, todos)
}

func (h *Handler) CreateTodo(w http.ResponseWriter, r *http.Request) {
	var req CreateTodoRequest

	if err := utils.DecodeJSON(r, &req); err != nil {
		slog.Warn(
			"invalid create todo request",
			"error", err,
		)

		utils.WriteJSONResponseError(w, http.StatusBadRequest, "invalid request body: "+err.Error())
		return
	}

	todo := items.TodoItem{
		Title: req.Title,
	}

	ctx := r.Context()

	createdTodo, err := h.service.CreateTodo(ctx, todo)
	if err != nil {
		slog.Error("failed to create todo", "error", err, "todo_title", todo.Title)
		utils.WriteJSONResponseError(w, http.StatusBadGateway, "create todo error: "+err.Error())
		return
	}

	utils.WriteJSONResponse(w, http.StatusCreated, createdTodo)
}

func (h *Handler) DeleteTodoByID(w http.ResponseWriter, r *http.Request) {
	rawID := r.PathValue("id")

	id, err := strconv.Atoi(rawID)
	if err != nil {
		slog.Warn("invalid todo id", "error", err, "id", rawID)
		utils.WriteJSONResponseError(w, http.StatusBadRequest, "invalid id: "+err.Error())
		return
	}

	ctx := r.Context()

	if err := h.service.DeleteTodoByID(ctx, id); err != nil {
		if errors.Is(err, customerrors.ErrNotFound) {
			slog.Debug("failed to delete todo by id not found", "error", err, "id", id)
			utils.WriteJSONResponseError(w, http.StatusNotFound, "not found")
			return
		}

		slog.Error("failed to delete todo by id", "error", err, "id", id)
		utils.WriteJSONResponseError(w, http.StatusConflict, "DeleteTodoByID error: "+err.Error())
		return
	}

	utils.WriteJSONResponse(w, http.StatusOK, DeleteTodoByIDResponse{
		OK: true,
	})
}

func (h *Handler) UpdateTodoByID(w http.ResponseWriter, r *http.Request) {
	var req UpdateTodoByIDRequest

	if err := utils.DecodeJSON(r, &req); err != nil {
		slog.Warn(
			"invalid create todo request",
			"error", err,
		)
		utils.WriteJSONResponseError(w, http.StatusBadRequest, "invalid request body: "+err.Error())
		return
	}

	todo := items.TodoItem{
		ID:    req.ID,
		Title: req.Title,
	}

	ctx := r.Context()

	newTodo, err := h.service.UpdateTodoByID(ctx, todo)
	if err != nil {
		if errors.Is(err, customerrors.ErrNotFound) {
			slog.Debug("failed to update todo by id todo not found", "error", err, "id", todo.ID, "title", todo.Title)
			utils.WriteJSONResponseError(w, http.StatusNotFound, "not found")
			return
		}

		slog.Error("failed to update todo by id", "error", err, "id", todo.ID, "title", todo.Title)
		utils.WriteJSONResponseError(w, http.StatusConflict, "UpdateTodoByID error: "+err.Error())
		return
	}

	utils.WriteJSONResponse(w, http.StatusOK, UpdateTodoByIDResponse{
		ID:    newTodo.ID,
		Title: newTodo.Title,
	})
}
