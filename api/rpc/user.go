package rpc

import (
	"context"
	"errors"
	"time"

	"github.com/zjutjh/User-Center/comm"
	"github.com/zjutjh/User-Center/comm/bizerr"
	"github.com/zjutjh/User-Center/comm/response"
	"github.com/zjutjh/mygo/nlog"

	"github.com/zjutjh/User-Center/dao"
	"github.com/zjutjh/User-Center/dao/model"
	"gorm.io/gorm"

	pb "github.com/zjutjh/User-Center/idl/user/v1alpha1"
	"github.com/zjutjh/WeJH-SDK/oauth"
	"github.com/zjutjh/WeJH-SDK/oauth/oauthException"
)

type UserService struct {
	pb.UnimplementedUserCenterServiceServer
}

func NewUserService() *UserService {
	return &UserService{}
}

func (u *UserService) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	if err := u.checkStudentBySIDAndIID(ctx, req.StudentId, req.Iid); err != nil {
		return nil, bizerr.UserNotExist
	}
	user, err := dao.GetUserByStudentId(ctx, req.StudentId)
	if user != nil {
		return &pb.RegisterResponse{}, nil
	}
	if err != nil {
		nlog.Pick().Errorf("failed to get user by student id: %v", err)
		return response.Error[pb.RegisterResponse](err)
	}
	pass := comm.Encrypt(req.Password)
	user = &model.User{
		Password:   pass,
		StudentID:  req.StudentId,
		Email:      req.Email,
		CreateTime: time.Now(),
	}
	if err := dao.CreateUser(ctx, user); err != nil {
		nlog.Pick().Errorf("failed to create user by student id: %v", err)
		return response.Error[pb.RegisterResponse](bizerr.UnknownError)
	}
	return response.OK[pb.RegisterResponse]()
}

func (u *UserService) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	user, err := dao.GetUserByStudentId(ctx, req.StudentId)
	if err != nil {
		return response.Error[pb.LoginResponse](bizerr.UserNotExist)
	}
	if user == nil {
		return response.Error[pb.LoginResponse](bizerr.UserNotExist)
	}
	if user.Password != comm.Encrypt(req.Password) {
		return response.Error[pb.LoginResponse](bizerr.WrongAccountOrPassword)
	}
	return response.OK[pb.LoginResponse]()
}

func (u *UserService) ResetPassword(ctx context.Context, req *pb.ResetPasswordRequest) (*pb.ResetPasswordResponse, error) {
	if err := u.checkStudentBySIDAndIID(ctx, req.StudentId, req.Iid); err != nil {
		return response.Error[pb.ResetPasswordResponse](bizerr.UserNotExist)
	}
	user, err := dao.GetUserByStudentId(ctx, req.StudentId)
	if err != nil {
		return response.Error[pb.ResetPasswordResponse](bizerr.UnknownError)
	}
	user.Password = comm.Encrypt(req.Password)
	if err := dao.UpdateUserPassword(ctx, user); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return response.Error[pb.ResetPasswordResponse](bizerr.UserNotExist)
		}
		nlog.Pick().Errorf("failed to update user password: %v", err)
		return response.Error[pb.ResetPasswordResponse](bizerr.UnknownError)
	}
	return response.OK[pb.ResetPasswordResponse]()
}

func (u *UserService) Delete(ctx context.Context, req *pb.DeleteRequest) (*pb.DeleteResponse, error) {
	if err := u.checkStudentBySIDAndIID(ctx, req.StudentId, req.Iid); err != nil {
		return response.Error[pb.DeleteResponse](bizerr.UserNotExist)
	}
	if err := dao.DeleteUserByStudentId(ctx, req.StudentId); err != nil {
		return response.Error[pb.DeleteResponse](bizerr.UnknownError)
	}
	return response.OK[pb.DeleteResponse]()
}

func (u *UserService) OauthLogin(ctx context.Context, req *pb.LoginRequest) (*pb.OauthLoginResponse, error) {
	_, userInfo, err := oauth.GetUserInfo(req.StudentId, req.Password)
	if err != nil {
		switch {
		case errors.Is(err, oauthException.ClosedError):
			return response.Error[pb.OauthLoginResponse](bizerr.ClosedError)
		case errors.Is(err, oauthException.WrongPassword):
			return response.Error[pb.OauthLoginResponse](bizerr.WrongAccountOrPassword)
		case errors.Is(err, oauthException.NotActivatedError):
			return response.Error[pb.OauthLoginResponse](bizerr.NotActivated)
		case errors.Is(err, oauthException.WrongAccount):
			return response.Error[pb.OauthLoginResponse](bizerr.WrongAccountOrPassword)
		case errors.Is(err, oauthException.OtherError):
			return response.Error[pb.OauthLoginResponse](bizerr.UnknownError)
		default:
			return response.Error[pb.OauthLoginResponse](bizerr.UnknownError)
		}
	}
	return &pb.OauthLoginResponse{
		Code:         pb.BizCode_OK,
		Name:         userInfo.Name,
		StudentId:    userInfo.StudentID,
		UserType:     userInfo.UserType,
		UserTypeDesc: userInfo.UserTypeDesc,
		Gender:       userInfo.Gender,
		Avatar:       userInfo.Avatar,
	}, nil
}

func (u *UserService) checkStudentBySIDAndIID(ctx context.Context, studentID string, iid string) error {
	student, err := dao.GetStudentByStudentID(ctx, studentID)
	if err != nil || student == nil || student.IDCard != iid {
		return bizerr.UserNotExist
	}
	return nil
}
