package repository

import (
	"authService/internal/domain"
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	pool *pgxpool.Pool
}

func NewUserRepository(pool *pgxpool.Pool) *UserRepository {
	return &UserRepository{pool}
}

func (ur *UserRepository) Exists(ctx context.Context, email string) (bool, error) {
	query := `
		SELECT id FROM users WHERE email = $1;
	`
	var id string
	err := ur.pool.QueryRow(ctx, query, email).Scan(&id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, nil
		}
		return true, err
	}
	return true, nil
}

func (ur *UserRepository) Save(ctx context.Context, user *domain.User) error {
	query := `
		INSERT INTO users (email, password, username)
		VALUES ($1, $2, $3)
		RETURNING id, created_at;
	`
	err := ur.pool.QueryRow(ctx, query, user.Email, user.Password, user.Name).Scan(&user.ID, &user.CreatedAt)
	if err != nil {
		return fmt.Errorf("error saving user: %w", err)
	}
	return nil
}
