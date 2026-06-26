package logic

import (
	"context"

	"github.com/zjutjh/User-Center/apps/user-rpc/internal/svc"
	"github.com/zjutjh/User-Center/apps/user-rpc/pb"
	"github.com/zjutjh/User-Center/common/errorsx"

	"github.com/zeromicro/go-zero/core/logx"
)

type BindZfLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewBindZfLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BindZfLogic {
	return &BindZfLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *BindZfLogic) BindZf(in *pb.BindZfRequest) (*pb.BindZfResponse, error) {
	zfPassword, err := l.svcCtx.Credential.PrepareZfBind(in.GetZfPassword())
	if err != nil {
		if codeErr, ok := errorsx.As(err); ok {
			return nil, codeErr
		}
		l.Errorf("构建用户正方绑定信息失败: %v", err)
		return nil, errorsx.ErrUnknown
	}

	if err := ensureBindUserExists(l.ctx, l.svcCtx, l.Logger, in.GetUserId()); err != nil {
		return nil, err
	}
	if err := l.svcCtx.UserRepo.UpdateZfBindByID(l.ctx, in.GetUserId(), zfPassword); err != nil {
		l.Errorf("更新用户正方绑定信息失败: %v", err)
		return nil, errorsx.ErrUnknown
	}

	return &pb.BindZfResponse{}, nil
}
