package logic

import (
	"context"
	"errors"

	"github.com/zjutjh/User-Center/apps/user-rpc/internal/dao/model"
	"github.com/zjutjh/User-Center/apps/user-rpc/internal/svc"
	"github.com/zjutjh/User-Center/apps/user-rpc/pb"
	"github.com/zjutjh/User-Center/common/errorsx"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type BindLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewBindLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BindLogic {
	return &BindLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *BindLogic) Bind(in *pb.BindRequest) (*pb.BindResponse, error) {
	if in.UserId <= 0 {
		return nil, errorsx.ErrParameterInvalid
	}

	_, err := l.svcCtx.Query.User.WithContext(l.ctx).
		Where(l.svcCtx.Query.User.ID.Eq(in.UserId)).
		First()
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errorsx.ErrUserNotExist
	}
	if err != nil {
		l.Errorf("根据用户 ID 查询用户失败: %v", err)
		return nil, errorsx.ErrUnknown
	}

	updates, err := l.svcCtx.Credential.BuildBindUpdates(in)
	if err != nil {
		if errors.Is(err, errorsx.ErrParameterInvalid) {
			return nil, errorsx.ErrParameterInvalid
		}
		l.Errorf("构建用户绑定信息失败: %v", err)
		return nil, errorsx.ErrUnknown
	}

	if err := l.svcCtx.DB.WithContext(l.ctx).
		Model(&model.User{}).
		Where("id = ?", in.UserId).
		Updates(updates).Error; err != nil {
		l.Errorf("更新用户绑定信息失败: %v", err)
		return nil, errorsx.ErrUnknown
	}

	return &pb.BindResponse{}, nil
}
