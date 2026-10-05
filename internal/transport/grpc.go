package transport

import (
	"authService/internal/service"
	"context"
	"fmt"

	authv1 "github.com/eduardkarpow/auth-proto/gen/go/auth/v1"
)

type UserGrpcHandler struct {
	authv1.UnimplementedAuthServer
	userService *service.UserService
}

func (h *UserGrpcHandler) Register(ctx context.Context, request *authv1.RegisterRequest) (*authv1.TokensResponse, error) {
	resp, err := h.userService.Register(ctx, request)
	fmt.Println(err)
	return resp, err
}

func (h *UserGrpcHandler) Login(ctx context.Context, request *authv1.LoginRequest) (*authv1.TokensResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (h *UserGrpcHandler) Refresh(ctx context.Context, request *authv1.RefreshRequest) (*authv1.TokensResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (h *UserGrpcHandler) Logout(ctx context.Context, request *authv1.LoginRequest) (*authv1.SuccessResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (h *UserGrpcHandler) ResetPassword(ctx context.Context, request *authv1.ResetPasswordRequest) (*authv1.SuccessResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (h *UserGrpcHandler) mustEmbedUnimplementedAuthServer() {
	//TODO implement me
	panic("implement me")
}

func NewUserGrpcHandler(userService *service.UserService) *UserGrpcHandler {
	return &UserGrpcHandler{userService: userService}
}
