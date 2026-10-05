package domain

import (
	"context"
	"time"

	authv1 "github.com/eduardkarpow/auth-proto/gen/go/auth/v1"
)

type User struct {
	ID        string
	Email     string
	Password  string
	Name      string
	CreatedAt time.Time
}

type UserRepositoryInterface interface {
	Exists(ctx context.Context, email string) (bool, error)
	Save(ctx context.Context, user *User) error
	FindByEmail(ctx context.Context, email string) (*User, error)
}

type UserServiceInterface interface {
	Register(ctx context.Context, user authv1.RegisterRequest) (*Tokens, error)
	Login(ctx context.Context, user authv1.LoginRequest) (*Tokens, error)
	generateTokens(ctx context.Context, user *User) (*authv1.TokensResponse, error)
}
