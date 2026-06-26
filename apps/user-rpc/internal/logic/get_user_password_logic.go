package logic

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zjutjh/User-Center/apps/user-rpc/internal/svc"
	"github.com/zjutjh/User-Center/apps/user-rpc/pb"
	"github.com/zjutjh/User-Center/common/errorsx"
)

type GetUserPasswordLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUserPasswordLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserPasswordLogic {
	return &GetUserPasswordLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetUserPasswordLogic) GetUserPassword(in *pb.GetUserPasswordRequest) (*pb.GetUserPasswordResponse, error) {
	user, err := l.svcCtx.UserRepo.GetUserById(l.ctx, in.UserId)
	if err != nil {
		if codeErr, ok := errorsx.As(err); ok {
			return nil, codeErr
		}
		l.Errorf("根据用户 ID 查询用户失败: %v", err)
		return nil, errorsx.ErrUnknown
	}

	zfPassword, err := l.svcCtx.Credential.DecryptPassword(user.ZfPassword)
	if err != nil {
		l.Errorf("解密正方密码失败: %v", err)
		return nil, errorsx.ErrUnknown
	}

	oauthPassword, err := l.svcCtx.Credential.DecryptPassword(user.OauthPassword)
	if err != nil {
		l.Errorf("解密统一认证密码失败: %v", err)
		return nil, errorsx.ErrUnknown
	}

	return &pb.GetUserPasswordResponse{
		StudentId:     user.StudentID,
		DeviceId:      user.DeviceID,
		YxyUid:        user.YxyUID,
		ZfPassword:    zfPassword,
		OauthPassword: oauthPassword,
	}, nil
}
