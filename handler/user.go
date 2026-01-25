package handler

import (
	"context"
	"errors"

	"gorm.io/gorm"

	userv1 "github.com/zjutjh/User-Center-grpc/api/user/v1alpha1"
	"github.com/zjutjh/User-Center-grpc/pkg/expection"
	"github.com/zjutjh/User-Center-grpc/pkg/util"
	userService "github.com/zjutjh/User-Center-grpc/service/user"
	"github.com/zjutjh/WeJH-SDK/oauth"
	"github.com/zjutjh/WeJH-SDK/oauth/oauthException"
)

type UserHandler struct {
	userv1.UnimplementedUserCenterServiceServer
}

func NewUserHandler() *UserHandler {
	return &UserHandler{}
}

func (u *UserHandler) Register(ctx context.Context, req *userv1.RegisterRequest) (*userv1.Response, error) {
	if err := userService.CheckStudentBySIDAndIID(ctx, req.StudentId, req.Iid); err != nil {
		return expection.UserNotFound.ToResponse()
	}
	if err := userService.CreateUser(ctx, req.Password, req.Email, req.StudentId); err != nil {
		if errors.Is(err, expection.UserAlreadyExit) {
			return expection.UserAlreadyExit.ToResponse()
		}
		return expection.Unknown.ToResponse()
	}
	return util.ResponseSuccess(nil)
}

func (u *UserHandler) Login(ctx context.Context, req *userv1.LoginRequest) (*userv1.Response, error) {
	_, err := userService.Login(ctx, req.StudentId, req.Password)
	if err != nil {
		return expection.Unknown.ToResponse()
	}
	return util.ResponseSuccess(nil)
}

func (u *UserHandler) ResetPassword(ctx context.Context, req *userv1.ResetPasswordRequest) (*userv1.Response, error) {
	if err := userService.CheckStudentBySIDAndIID(ctx, req.StudentId, req.Iid); err != nil {
		return expection.UserNotFound.ToResponse()
	}
	if err := userService.UpdateUserPassword(ctx, req.StudentId, req.Password); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return expection.UserNotExit.ToResponse()
		}
		return expection.Unknown.ToResponse()
	}
	return util.ResponseSuccess(nil)
}

func (u *UserHandler) Delete(ctx context.Context, req *userv1.DeleteRequest) (*userv1.Response, error) {
	if err := userService.CheckStudentBySIDAndIID(ctx, req.StudentId, req.Iid); err != nil {
		return expection.UserNotFound.ToResponse()
	}
	if err := userService.Delete(ctx, req.StudentId); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return expection.UserNotExit.ToResponse()
		}
		return expection.Unknown.ToResponse()
	}
	return util.ResponseSuccess(nil)
}

func (u *UserHandler) OauthLogin(ctx context.Context, req *userv1.LoginRequest) (*userv1.Response, error) {
	_, userInfo, err := oauth.GetUserInfo(req.StudentId, req.Password)
	if err != nil {
		switch {
		case errors.Is(err, oauthException.ClosedError):
			return expection.ClosedError.ToResponse()
		case errors.Is(err, oauthException.WrongPassword):
			return expection.WrongPassword.ToResponse()
		case errors.Is(err, oauthException.NotActivatedError):
			return expection.NotActivatedError.ToResponse()
		case errors.Is(err, oauthException.WrongAccount):
			return expection.WrongAccount.ToResponse()
		case errors.Is(err, oauthException.OtherError):
			return expection.Unknown.ToResponse()
		}
	}
	return util.ResponseSuccess(map[string]interface{}{
		"name":         userInfo.Name,
		"studentId":    userInfo.StudentID,
		"userType":     userInfo.UserType,
		"userTypeDesc": userInfo.UserTypeDesc,
		"gender":       userInfo.Gender,
		"avatar":       userInfo.Avatar,
	})
}
