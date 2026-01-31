package userService

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/zjutjh/User-Center/biz/dao"
	"github.com/zjutjh/User-Center/biz/dao/model"
	"github.com/zjutjh/User-Center/biz/util/response"
	"github.com/zjutjh/User-Center/biz/util/security"
	"gorm.io/gorm"

	userv1 "github.com/zjutjh/User-Center/api/user/v1alpha1"
	"github.com/zjutjh/WeJH-SDK/oauth"
	"github.com/zjutjh/WeJH-SDK/oauth/oauthException"

	"github.com/zjutjh/User-Center/biz/err"
)

type UserService struct {
}

func NewUserService() *UserService {
	return &UserService{}
}

func (u *UserService) Register(ctx context.Context, req *userv1.RegisterRequest) (*userv1.Response, error) {
	if err := u.checkStudentBySIDAndIID(ctx, req.StudentId, req.Iid); err != nil {
		return nil, bizerr.UserNotExist
	}
	user, err := dao.GetUserByStudentId(ctx, req.StudentId)
	if user != nil {
		return response.Error(bizerr.UserExisted)
	}
	if err != nil {
		slog.Error("failed to get user by student id: %v", err)
		return response.Error(bizerr.UnknownError)
	}
	pass := security.Encrypt(req.Password)
	user = &model.User{
		Password:   pass,
		StudentID:  req.StudentId,
		Email:      req.Email,
		CreateTime: time.Now(),
	}
	if err := dao.CreateUser(ctx, user); err != nil {
		slog.Error("failed to create user by student id: %v", err)
		return response.Error(bizerr.UnknownError)
	}
	return response.Success(nil)
}

func (u *UserService) Login(ctx context.Context, req *userv1.LoginRequest) (*userv1.Response, error) {
	user, err := dao.GetUserByStudentId(ctx, req.StudentId)
	if err != nil {
		return response.Error(bizerr.UnknownError)
	}
	if user == nil {
		return response.Error(bizerr.UserNotExist)
	}
	if user.Password != security.Encrypt(req.Password) {
		return response.Error(bizerr.WrongAccountOrPassword)
	}
	return response.Success(nil)
}

func (u *UserService) ResetPassword(ctx context.Context, req *userv1.ResetPasswordRequest) (*userv1.Response, error) {
	if err := u.checkStudentBySIDAndIID(ctx, req.StudentId, req.Iid); err != nil {
		return response.Error(bizerr.UserNotExist)
	}
	user, err := dao.GetUserByStudentId(ctx, req.StudentId)
	if err != nil {
		return response.Error(bizerr.UnknownError)
	}
	user.Password = security.Encrypt(req.Password)
	if err := dao.UpdateUserPassword(ctx, user); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return response.Error(bizerr.UserNotExist)
		}
		slog.Error("failed to update user password: %v", err)
		return response.Error(bizerr.UnknownError)
	}
	return response.Success(nil)
}

func (u *UserService) Delete(ctx context.Context, req *userv1.DeleteRequest) (*userv1.Response, error) {
	if err := u.checkStudentBySIDAndIID(ctx, req.StudentId, req.Iid); err != nil {
		return response.Error(bizerr.UserNotExist)
	}
	if err := dao.DeleteUserByStudentId(ctx, req.StudentId); err != nil {
		return response.Error(bizerr.UnknownError)
	}
	return response.OK()
}

func (u *UserService) OauthLogin(ctx context.Context, req *userv1.LoginRequest) (*userv1.Response, error) {
	_, userInfo, err := oauth.GetUserInfo(req.StudentId, req.Password)
	if err != nil {
		switch {
		case errors.Is(err, oauthException.ClosedError):
			return response.Error(bizerr.ClosedError)
		case errors.Is(err, oauthException.WrongPassword):
			return response.Error(bizerr.WrongAccountOrPassword)
		case errors.Is(err, oauthException.NotActivatedError):
			return response.Error(bizerr.NotActivated)
		case errors.Is(err, oauthException.WrongAccount):
			return response.Error(bizerr.WrongAccountOrPassword)
		case errors.Is(err, oauthException.OtherError):
			return response.Error(bizerr.UnknownError)
		}
	}
	return response.Success(map[string]interface{}{
		"name":         userInfo.Name,
		"studentId":    userInfo.StudentID,
		"userType":     userInfo.UserType,
		"userTypeDesc": userInfo.UserTypeDesc,
		"gender":       userInfo.Gender,
		"avatar":       userInfo.Avatar,
	})
}

func (u *UserService) checkStudentBySIDAndIID(ctx context.Context, studentID string, iid string) error {
	student, err := dao.GetStudentByStudentID(ctx, studentID)
	if err != nil || student == nil || student.IDCard != iid {
		return bizerr.UserNotExist
	}
	return nil
}
