package handler

import (
	"context"

	userv1 "github.com/zjutjh/User-Center/api/user/v1alpha1"
	userService "github.com/zjutjh/User-Center/service/user"
)

type UserHandler struct {
	userv1.UnimplementedUserCenterServiceServer
	userService *userService.UserService
}

func NewUserHandler() *UserHandler {
	return &UserHandler{
		userService: userService.NewUserService(),
	}
}

func (h *UserHandler) Register(ctx context.Context, req *userv1.RegisterRequest) (*userv1.Response, error) {
	return h.userService.Register(ctx, req)
}

func (h *UserHandler) Login(ctx context.Context, req *userv1.LoginRequest) (*userv1.Response, error) {
	return h.userService.Login(ctx, req)
}

func (h *UserHandler) ResetPassword(ctx context.Context, req *userv1.ResetPasswordRequest) (*userv1.Response, error) {
	return h.userService.ResetPassword(ctx, req)
}

func (h *UserHandler) Delete(ctx context.Context, req *userv1.DeleteRequest) (*userv1.Response, error) {
	return h.userService.Delete(ctx, req)
}

func (h *UserHandler) OauthLogin(ctx context.Context, req *userv1.LoginRequest) (*userv1.Response, error) {
	return h.userService.OauthLogin(ctx, req)
}
