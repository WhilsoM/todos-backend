package repository

import (
	"context"
	"todos/internal/items"

	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{
		db,
	}
}

func (r *UserRepository) CreateUser(ctx context.Context, user items.User) (items.User, error) {
	query := `INSERT INTO users (email, password_hash) VALUES ($1, $2) RETURNING id, email`

	var newUser items.User

	row := r.db.QueryRow(ctx, query, user.Email, user.PasswordHash)

	if err := row.Scan(&newUser.ID, &newUser.Email); err != nil {
		return items.User{}, err
	}

	return newUser, nil
}

func (r *UserRepository) GetUserByEmail(ctx context.Context, email string) (items.User, error) {
	query := `SELECT id, email, password_hash FROM users WHERE email = $1`

	var user items.User

	row := r.db.QueryRow(ctx, query, email)

	if err := row.Scan(&user.ID, &user.Email, &user.PasswordHash); err != nil {
		return items.User{}, err
	}

	return user, nil
}
