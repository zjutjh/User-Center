package logic

import (
	"context"
	"errors"
	"strings"

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

type bindStrategy func(in *pb.BindRequest) (map[string]any, error)

var bindStrategies = map[pb.BindType]bindStrategy{
	pb.BindType_BIND_TYPE_YXY:   buildYXYBindUpdates,
	pb.BindType_BIND_TYPE_ZF:    buildZFBindUpdates,
	pb.BindType_BIND_TYPE_OAUTH: buildOAuthBindUpdates,
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

	updates, err := buildBindUpdates(in)
	if err != nil {
		return nil, err
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

func buildBindUpdates(in *pb.BindRequest) (map[string]any, error) {
	strategy, ok := bindStrategies[in.Type]
	if !ok {
		return nil, errorsx.ErrParameterInvalid
	}

	return strategy(in)
}

func buildYXYBindUpdates(in *pb.BindRequest) (map[string]any, error) {
	deviceID := strings.TrimSpace(in.DeviceId)
	yxyUID := strings.TrimSpace(in.YxyUid)
	if deviceID == "" || yxyUID == "" {
		return nil, errorsx.ErrParameterInvalid
	}

	return map[string]any{
		"device_id": deviceID,
		"yxy_uid":   yxyUID,
	}, nil
}

func buildZFBindUpdates(in *pb.BindRequest) (map[string]any, error) {
	zfPassword := strings.TrimSpace(in.ZfPassword)
	if zfPassword == "" {
		return nil, errorsx.ErrParameterInvalid
	}

	return map[string]any{
		"zf_password": zfPassword,
	}, nil
}

func buildOAuthBindUpdates(in *pb.BindRequest) (map[string]any, error) {
	oauthPassword := strings.TrimSpace(in.OauthPassword)
	if oauthPassword == "" {
		return nil, errorsx.ErrParameterInvalid
	}

	return map[string]any{
		"oauth_password": oauthPassword,
	}, nil
}
