package userService

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"log/slog"
	"time"

	userv1 "github.com/zjutjh/User-Center/api/user/v1alpha1"
	"github.com/zjutjh/WeJH-SDK/oauth"
	"github.com/zjutjh/WeJH-SDK/oauth/oauthException"

	"github.com/zjutjh/User-Center/dao"
	"github.com/zjutjh/User-Center/dao/model"
	"github.com/zjutjh/User-Center/pkg/expection"
	"github.com/zjutjh/User-Center/pkg/util"
)

type UserService struct {
}

func NewUserService() *UserService {
	return &UserService{}
}

func (u *UserService) Register(ctx context.Context, req *userv1.RegisterRequest) (*userv1.Response, error) {
	if err := u.checkStudentBySIDAndIID(ctx, req.StudentId, req.Iid); err != nil {
		return expection.UserNotFound.ToResponse()
	}
	user, err := dao.GetUserByStudentId(ctx, req.StudentId)
	if user != nil {
		return expection.UserAlreadyExit.ToResponse()
	}
	if err != nil {
		slog.Error("failed to get user by student id: %v", err)
		return expection.Unknown.ToResponse()
	}
	pass := util.Encryrpt(req.Password)
	user = &model.User{
		Password:   pass,
		StudentID:  req.StudentId,
		Email:      req.Email,
		CreateTime: time.Now(),
	}
	if err := dao.CreateUser(ctx, user); err != nil {
		slog.Error("failed to create user by student id: %v", err)
		return expection.Unknown.ToResponse()
	}
	return util.ResponseSuccess(nil)
}

func (u *UserService) Login(ctx context.Context, req *userv1.LoginRequest) (*userv1.Response, error) {
	user, err := dao.GetUserByStudentId(ctx, req.StudentId)
	if err != nil {
		return nil, expection.Unknown
	}
	if user == nil {
		return nil, expection.UserNotExit
	}
	if user.Password != util.Encryrpt(req.Password) {
		return nil, expection.AuthError
	}
	return util.ResponseSuccess(nil)
}

func (u *UserService) ResetPassword(ctx context.Context, req *userv1.ResetPasswordRequest) (*userv1.Response, error) {
	if err := u.checkStudentBySIDAndIID(ctx, req.StudentId, req.Iid); err != nil {
		return expection.UserNotFound.ToResponse()
	}
	user, err := dao.GetUserByStudentId(ctx, req.StudentId)
	if err != nil {
		return expection.Unknown.ToResponse()
	}
	user.Password = util.Encryrpt(req.Password)
	if err := dao.UpdateUserPassword(ctx, user); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return expection.UserNotExit.ToResponse()
		}
		slog.Error("failed to update user password: %v", err)
		return expection.Unknown.ToResponse()
	}
	return util.ResponseSuccess(nil)
}

func (u *UserService) Delete(ctx context.Context, req *userv1.DeleteRequest) (*userv1.Response, error) {
	if err := u.checkStudentBySIDAndIID(ctx, req.StudentId, req.Iid); err != nil {
		return expection.UserNotFound.ToResponse()
	}
	if err := dao.DeleteUserByStudentId(ctx, req.StudentId); err != nil {
		return expection.Unknown.ToResponse()
	}
	return util.ResponseSuccess(nil)
}

func (u *UserService) OauthLogin(ctx context.Context, req *userv1.LoginRequest) (*userv1.Response, error) {
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

func (u *UserService) checkStudentBySIDAndIID(ctx context.Context, studentID string, iid string) error {
	student, err := dao.GetStudentByStudentID(ctx, studentID)
	if err != nil || student == nil || student.IDCard != iid {
		return expection.UserNotFound
	}
	return nil
}

func (u *UserService) HealthyCheck(ctx context.Context, req *userv1.HealthyCheckRequest) (*userv1.Response, error) {
	return util.ResponseSuccess(map[string]interface{}{
		"studentID": req.StudentId,
	})
}
