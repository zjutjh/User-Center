package logic

import (
	"context"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zjutjh/User-Center/apps/user-rpc/internal/svc"
	"github.com/zjutjh/User-Center/apps/user-rpc/pb"
	"github.com/zjutjh/User-Center/common/errorsx"
)

type GetUserInfoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUserInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserInfoLogic {
	return &GetUserInfoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetUserInfoLogic) GetUserInfo(in *pb.GetUserInfoRequest) (*pb.GetUserInfoResponse, error) {
	if in.UserId <= 0 {
		return nil, errorsx.ErrParameterInvalid
	}

	user, err := getUserByID(l.ctx, l.svcCtx, in.UserId)
	if err != nil {
		if codeErr, ok := errorsx.As(err); ok {
			return nil, codeErr
		}
		l.Errorf("根据用户 ID 查询用户信息失败: %v", err)
		return nil, errorsx.ErrUnknown
	}

	return &pb.GetUserInfoResponse{
		UserId:        user.ID,
		StudentId:     user.StudentID,
		UserType:      user.Type,
		Email:         user.Email,
		PhoneNum:      user.PhoneNum,
		DeviceId:      user.DeviceID,
		YxyUid:        user.YxyUID,
		ZfPassword:    user.ZfPassword,
		OauthPassword: user.OauthPassword,
		CreateTime:    user.CreateTime.Format(time.RFC3339),
	}, nil
}
