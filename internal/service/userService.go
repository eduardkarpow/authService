package service

import (
	"authService/internal/config"
	"authService/internal/domain"
	"authService/internal/repository"
	"context"
	"errors"

	authv1 "github.com/eduardkarpow/auth-proto/gen/go/auth/v1"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserService struct {
	r *repository.UserRepository
	t *repository.TokenRepository
	s *SecurityService
}

func NewUserService(pool *pgxpool.Pool, config *config.Config) *UserService {
	return &UserService{
		r: repository.NewUserRepository(pool),
		t: repository.NewTokenRepository(pool),
		s: NewSecurityService(config.JWTSecret),
	}
}

func (us *UserService) Register(ctx context.Context, user *authv1.RegisterRequest) (*authv1.TokensResponse, error) {
	exists, err := us.r.Exists(ctx, user.Email)
	if err != nil {
		return &authv1.TokensResponse{}, err
	}
	if exists {
		return &authv1.TokensResponse{}, errors.New("user already exists")
	}
	passwordHashed, err := us.s.HashPassword(user.Password)
	if err != nil {
		return &authv1.TokensResponse{}, err
	}
	userInner := &domain.User{
		Email:    user.Email,
		Password: passwordHashed,
		Name:     user.Name,
	}
	err = us.r.Save(ctx, userInner)
	if err != nil {
		return &authv1.TokensResponse{}, err
	}
	tokens, err := us.s.GenerateTokens(userInner.ID)
	if err != nil {
		return &authv1.TokensResponse{}, err
	}
	tokenHash, err := us.s.HashPassword(tokens.RefreshToken)
	if err != nil {
		return &authv1.TokensResponse{}, err
	}
	session := &domain.Session{
		UserId:    userInner.ID,
		TokenHash: tokenHash,
		ExpiresAt: tokens.ExpiresIn,
	}
	err = us.t.Save(ctx, session)
	if err != nil {
		return &authv1.TokensResponse{}, err
	}
	return &authv1.TokensResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
		ExpiresIn:    tokens.ExpiresIn.Unix(),
	}, nil
}
