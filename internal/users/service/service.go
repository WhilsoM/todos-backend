package service

import (
	"context"
	"errors"
	"log/slog"
	"strings"
	"time"
	"todos/internal/items"
	"todos/internal/jwt"

	"golang.org/x/crypto/bcrypt"
)

type UserEventProducer interface {
	PublishUserRegistered(ctx context.Context, userID int) error
}

type UserRepository interface {
	CreateUser(ctx context.Context, user items.User) (items.User, error)
	GetUserByEmail(ctx context.Context, email string) (items.User, error)
}

type UserService struct {
	repo       UserRepository
	jwtmanager *jwt.JWTManager

	eventProducer UserEventProducer
}

func NewUserService(repo UserRepository, jwtmanager *jwt.JWTManager, eventProducer UserEventProducer) *UserService {
	return &UserService{
		repo,
		jwtmanager,
		eventProducer,
	}
}

func (s *UserService) RegisterUser(ctx context.Context, email, password string) (items.User, error) {
	emailTrimmed := strings.TrimSpace(email)

	if emailTrimmed == "" || strings.TrimSpace(password) == "" {
		return items.User{}, errors.New("email or password cannot be empty")
	}
	start := time.Now()

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return items.User{}, err
	}
	slog.Info("bcrypt", "duration", time.Since(start))

	start = time.Now()
	userRepo := items.User{
		Email:        emailTrimmed,
		PasswordHash: string(hash),
	}

	newUser, err := s.repo.CreateUser(ctx, userRepo)
	if err != nil {
		return items.User{}, err
	}
	slog.Info("database", "duration", time.Since(start))

	start = time.Now()
	if err := s.eventProducer.PublishUserRegistered(ctx, newUser.ID); err != nil {
		return items.User{}, err
	}
	slog.Info("kafka", "duration", time.Since(start))

	return newUser, nil
}

func (s *UserService) LoginUser(ctx context.Context, email, password string) (string, error) {
	emailTrimmed := strings.TrimSpace(email)

	if emailTrimmed == "" || strings.TrimSpace(password) == "" {
		return "", errors.New("email or password cannot be empty")
	}

	user, err := s.repo.GetUserByEmail(ctx, emailTrimmed)
	if err != nil {
		return "", err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return "", err
	}

	token, err := s.jwtmanager.GenerateToken(user.ID)
	if err != nil {
		return "", err
	}

	return token, nil
}
