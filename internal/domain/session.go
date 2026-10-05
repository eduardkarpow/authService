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
	GenerateAccessToken(userId string) (string, error)
	GenerateRefreshToken(userId string) (string, error)
}
