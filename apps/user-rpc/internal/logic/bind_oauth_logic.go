package logic

import (
	"context"

	"github.com/zjutjh/User-Center/apps/user-rpc/internal/svc"
	"github.com/zjutjh/User-Center/apps/user-rpc/pb"
	"github.com/zjutjh/User-Center/common/errorsx"

	"github.com/zeromicro/go-zero/core/logx"
)

type BindOauthLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewBindOauthLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BindOauthLogic {
	return &BindOauthLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *BindOauthLogic) BindOauth(in *pb.BindOauthRequest) (*pb.BindOauthResponse, error) {
	oauthPassword, err := l.svcCtx.Credential.PrepareOauthBind(in.GetOauthPassword())
	if err != nil {
		if codeErr, ok := errorsx.As(err); ok {
			return nil, codeErr
		}
		l.Errorf("构建用户统一认证绑定信息失败: %v", err)
		return nil, errorsx.ErrUnknown
	}

	if err := ensureBindUserExists(l.ctx, l.svcCtx, l.Logger, in.GetUserId()); err != nil {
		return nil, err
	}
	if err := l.svcCtx.UserRepo.UpdateOauthBindByID(l.ctx, in.GetUserId(), oauthPassword); err != nil {
		l.Errorf("更新用户统一认证绑定信息失败: %v", err)
		return nil, errorsx.ErrUnknown
	}

	return &pb.BindOauthResponse{}, nil
}
