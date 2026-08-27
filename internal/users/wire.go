package users

import (
	"net/http"
	"todos/internal/jwt"
	"todos/internal/users/handler"
	"todos/internal/users/repository"
	"todos/internal/users/service"

	"github.com/jackc/pgx/v5/pgxpool"
)

func Register(router *http.ServeMux, db *pgxpool.Pool, jwtmanager *jwt.JWTManager, eventProducer service.UserEventProducer) {
	repo := repository.NewUserRepository(db)
	svc := service.NewUserService(repo, jwtmanager, eventProducer)
	h := handler.NewUserHandler(svc)

	router.HandleFunc("POST /api/register", h.RegisterUser)
	router.HandleFunc("POST /api/login", h.LoginUser)
}
