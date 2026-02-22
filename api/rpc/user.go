package rpc

import (
	"context"

	"github.com/zjutjh/User-Center/comm"
	"github.com/zjutjh/User-Center/comm/response"
	"github.com/zjutjh/User-Center/dao/model"
	"github.com/zjutjh/User-Center/dao/repo"
	pb "github.com/zjutjh/User-Center/idl/user/v1alpha1"
	"github.com/zjutjh/mygo/nlog"
)

type UserService struct {
	pb.UnimplementedUserCenterServiceServer
}

func (u *UserService) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	user, err := repo.NewUserRepo().GetUserByStudentId(ctx, req.StudentId)
	if err != nil {
		return response.Error[pb.LoginResponse](err)
	}
	if user == nil {
		return response.Error[pb.LoginResponse](comm.UserNotExist)
	}
	if user.Password != comm.Sha256Hash(req.Password) {
		return response.Error[pb.LoginResponse](comm.WrongAccountOrPassword)
	}
	return &pb.LoginResponse{Code: pb.BizCode_OK, UserId: user.ID}, nil
}

func (u *UserService) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	if _, err := repo.NewStudentRepo().CheckStudentByStudentIDAndCardId(ctx, req.StudentId, req.CardId); err != nil {
		return nil, comm.UserNotExist
	}
	newUser := &model.User{
		Password:  comm.Sha256Hash(req.Password),
		StudentID: req.StudentId,
		Email:     req.Email,
	}
	if err := repo.NewUserRepo().CreateUser(ctx, newUser); err != nil {
		nlog.Pick().WithError(err).Errorf("创建用户失败: %v", err)
		return response.Error[pb.RegisterResponse](comm.UnknownError)
	}
	return response.OK[pb.RegisterResponse]()
}

func (u *UserService) ResetPassword(ctx context.Context, request *pb.ResetPasswordRequest) (*pb.ResetPasswordResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (u *UserService) Delete(ctx context.Context, request *pb.DeleteRequest) (*pb.DeleteResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (u *UserService) OauthLogin(ctx context.Context, request *pb.LoginRequest) (*pb.OauthLoginResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (u *UserService) HealthyCheck(ctx context.Context, request *pb.HealthyCheckRequest) (*pb.HealthyCheckResponse, error) {
	//TODO implement me
	panic("implement me")
}
