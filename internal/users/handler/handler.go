package handler

import (
	"context"
	"log/slog"
	"net/http"
	customerrors "todos/internal/custom-errors"
	"todos/internal/items"
	"todos/internal/utils"
)

type UserService interface {
	RegisterUser(ctx context.Context, email, password string) (items.User, error)
	LoginUser(ctx context.Context, email, password string) (string, error)
}

type UserHandler struct {
	service UserService
}

func NewUserHandler(service UserService) *UserHandler {
	return &UserHandler{
		service,
	}
}

func (h *UserHandler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	var req items.RegisterUserRequest

	if err := utils.DecodeJSON(r, &req); err != nil {
		slog.Info("failed to decode json", "error", err, "email", req.Email)
		utils.WriteJSONResponseError(w, http.StatusBadRequest, customerrors.ErrBadRequest.Error())
		return
	}

	user, err := h.service.RegisterUser(r.Context(), req.Email, req.Password)
	if err != nil {
		slog.Info("failed to register user", "error", err)
		utils.WriteJSONResponseError(w, http.StatusBadGateway, "error: "+err.Error())
		return
	}

	res := items.RegisterUserResponse{
		ID:    user.ID,
		Email: user.Email,
	}
	utils.WriteJSONResponse(w, http.StatusCreated, res)
}

func (h *UserHandler) LoginUser(w http.ResponseWriter, r *http.Request) {
	var req items.LoginUserRequest

	if err := utils.DecodeJSON(r, &req); err != nil {
		slog.Info("failed to decode json", "error", err, "email", req.Email)
		utils.WriteJSONResponseError(w, http.StatusBadRequest, customerrors.ErrBadRequest.Error())
		return
	}

	token, err := h.service.LoginUser(r.Context(), req.Email, req.Password)
	if err != nil {
		slog.Info("failed to login user", "error", err, "email", req.Email)
		utils.WriteJSONResponseError(w, http.StatusBadGateway, err.Error())
		return
	}

	res := items.LoginUserResponse{
		Token: token,
	}

	utils.WriteJSONResponse(w, http.StatusOK, res)
}
