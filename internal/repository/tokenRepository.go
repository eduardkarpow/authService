package repository

import (
	"authService/internal/domain"
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type TokenRepository struct {
	pool *pgxpool.Pool
}

func NewTokenRepository(pool *pgxpool.Pool) *TokenRepository {
	return &TokenRepository{pool}
}

func (tr *TokenRepository) Save(ctx context.Context, session *domain.Session) error {
	query := `INSERT INTO sessions (user_id, token_hash, expires_at)
			  VALUES ($1, $2, $3)
	`
	_, err := tr.pool.Exec(ctx, query, session.UserId, session.TokenHash, session.ExpiresAt)
	return err
}
