package domain

import (
	"context"
	"time"
)

type Session struct {
	ID        string
	UserId    string
	TokenHash string
	UserAgent string
	ClientIp  string
	ExpiresAt time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Tokens struct {
	AccessToken  string
	RefreshToken string
	ExpiresIn    time.Time
}

type TokenRepositoryInterface interface {
	Save(ctx context.Context, session *Session) error
}

type SecurityServiceInterface interface {
	HashPassword(password string) (string, error)
	GenerateTokens(userId string) (Tokens, error)
	generateToken(userId string, expiry int64) (string, error)
	VerifyPassword(password string) error
	getSha256(password string) (string, error)
}
