package bff

import (
	"context"
	"errors"

	"gorm.io/gorm"

	typev1 "github.com/zjutjh/User-Center-grpc/api/types/v1alpha1"
	userv1 "github.com/zjutjh/User-Center-grpc/api/user/v1alpha1"
	apiv1 "github.com/zjutjh/User-Center-grpc/api/v1"
	"github.com/zjutjh/User-Center-grpc/pkg/apiExpection"
	"github.com/zjutjh/User-Center-grpc/pkg/services/user"
	"github.com/zjutjh/User-Center-grpc/pkg/util"
	"github.com/zjutjh/WeJH-SDK/oauth"
	"github.com/zjutjh/WeJH-SDK/oauth/oauthException"
)

type UserHandler struct {
	apiv1.UnimplementedUserServer
}

func NewUserHandler() *UserHandler {
	return &UserHandler{}
}

func (u *UserHandler) Register(ctx context.Context, req *userv1.RegisterRequest) (*typev1.Response, error) {
	if err := userService.CheckStudentBySIDAndIID(req.StudentId, req.Iid); err != nil {
		return apiExpection.UserNotFound.ToResponse()
	}
	if err := userService.CreateUser(req.Password, req.Email, req.StudentId); err != nil {
		if errors.Is(err, apiExpection.UserAlreadyExit) {
			return apiExpection.UserAlreadyExit.ToResponse()
		}
		return apiExpection.Unknown.ToResponse()
	}
	return util.ResponseSuccess(nil)
}

func (u *UserHandler) Login(ctx context.Context, req *userv1.LoginRequest) (*typev1.Response, error) {
	user, err := userService.GetUserByStudentId(req.StudentId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return apiExpection.UserNotExit.ToResponse()
		}
		return apiExpection.Unknown.ToResponse()
	}
	if user.Password != util.Encryrpt(req.Password) {
		return apiExpection.AuthError.ToResponse()
	}
	return util.ResponseSuccess(nil)
}

func (u *UserHandler) ResetPassword(ctx context.Context, req *userv1.ResetPasswordRequest) (*typev1.Response, error) {
	if err := userService.CheckStudentBySIDAndIID(req.StudentId, req.Iid); err != nil {
		return apiExpection.UserNotFound.ToResponse()
	}
	if err := userService.UpdateUserPassword(req.StudentId, req.Password); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return apiExpection.UserNotExit.ToResponse()
		}
		return apiExpection.Unknown.ToResponse()
	}
	return util.ResponseSuccess(nil)
}

func (u *UserHandler) Delete(ctx context.Context, req *userv1.DeleteRequest) (*typev1.Response, error) {
	if err := userService.CheckStudentBySIDAndIID(req.StudentId, req.Iid); err != nil {
		return apiExpection.UserNotFound.ToResponse()
	}
	if err := userService.Delete(req.StudentId); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return apiExpection.UserNotExit.ToResponse()
		}
		return apiExpection.Unknown.ToResponse()
	}
	return util.ResponseSuccess(nil)
}

func (u *UserHandler) OauthLogin(ctx context.Context, req *userv1.LoginRequest) (*typev1.Response, error) {
	_, err := oauth.Login(req.StudentId, req.Password)
	if err != nil {
		switch {
		case errors.Is(err, oauthException.ClosedError):
			return apiExpection.ClosedError.ToResponse()
		case errors.Is(err, oauthException.WrongPassword):
			return apiExpection.WrongPassword.ToResponse()
		case errors.Is(err, oauthException.NotActivatedError):
			return apiExpection.NotActivatedError.ToResponse()
		case errors.Is(err, oauthException.WrongAccount):
			return apiExpection.WrongAccount.ToResponse()
		case errors.Is(err, oauthException.OtherError):
			return apiExpection.Unknown.ToResponse()
		}
	}
	return util.ResponseSuccess(nil)
}
