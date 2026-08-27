package todos

import (
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"

	"todos/internal/jwt"
	"todos/internal/middleware"
	"todos/internal/todos/handler"
	cache "todos/internal/todos/redis"
	"todos/internal/todos/repository"
	"todos/internal/todos/service"
)

func Register(router *http.ServeMux, db *pgxpool.Pool, jwt *jwt.JWTManager, redisClient *redis.Client) {
	cache := cache.NewTodoCache(redisClient)

	repo := repository.NewRepository(db)
	svc := service.NewService(repo, cache)
	h := handler.NewHandler(svc)

	auth := middleware.AuthMiddleware(jwt)

	router.Handle("GET /api/todos", auth(http.HandlerFunc(h.GetTodos)))
	router.Handle("POST /api/todos", auth(http.HandlerFunc(h.CreateTodo)))
	router.Handle("DELETE /api/todos/{id}", auth(http.HandlerFunc(h.DeleteTodoByID)))
	router.Handle("PUT /api/todos", auth(http.HandlerFunc(h.UpdateTodoByID)))
}
