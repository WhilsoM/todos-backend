package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"todos/internal/config"
	"todos/internal/connections"
	"todos/internal/jwt"
	"todos/internal/todos"
	"todos/internal/users"
	"todos/internal/users/events"
)

func main() {
	cfg := config.MustConfig()
	router := http.NewServeMux()
	ctx := context.Background()

	jwtmanager := jwt.NewJWTManager(cfg.JWTSecret)

	userEventProducer := events.NewUserEventProducer(cfg.UserEventKafkaPort)
	defer userEventProducer.Close()

	db, err := connections.NewPool(ctx, cfg.DatabaseURL)
	if err != nil {
		slog.Error("failed to connect to database", "error", err)
		return
	}
	defer db.Close()

	redisClient, err := connections.NewRedisClient(ctx, cfg.RedisPort)
	if err != nil {
		slog.Error("failed to connect to redis", "error", err)
		return
	}
	defer redisClient.Close()

	users.Register(router, db, jwtmanager, userEventProducer)
	todos.Register(router, db, jwtmanager, redisClient)

	slog.Info("Server started work on port", "port", cfg.Port)

	if err := http.ListenAndServe(cfg.Port, router); err != nil {
		slog.Error("server stopped to work", "error", err)
		os.Exit(1)
	}
}
