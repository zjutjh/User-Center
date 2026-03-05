package logic

import (
	"context"
	"errors"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zjutjh/User-Center/apps/user-rpc/internal/svc"
	"github.com/zjutjh/User-Center/apps/user-rpc/pb"
	"github.com/zjutjh/User-Center/common/errorsx"
	"gorm.io/gorm"
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
	user, err := l.svcCtx.Query.User.WithContext(l.ctx).
		Where(l.svcCtx.Query.User.ID.Eq(in.UserId)).
		First()
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errorsx.ErrUserNotExist
	}
	if err != nil {
		l.Errorf("根据用户 ID 查询用户失败: %v", err)
		return nil, errorsx.ErrUnknown
	}

	return &pb.GetUserPasswordResponse{
		StudentId:     user.StudentID,
		DeviceId:      user.DeviceID,
		YxyUid:        user.YxyUID,
		ZfPassword:    user.ZfPassword,
		OauthPassword: user.OauthPassword,
	}, nil
}
