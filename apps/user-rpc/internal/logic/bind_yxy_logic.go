package logic

import (
	"context"

	"github.com/zjutjh/User-Center/apps/user-rpc/internal/svc"
	"github.com/zjutjh/User-Center/apps/user-rpc/pb"
	"github.com/zjutjh/User-Center/common/errorsx"

	"github.com/zeromicro/go-zero/core/logx"
)

type BindYxyLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewBindYxyLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BindYxyLogic {
	return &BindYxyLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *BindYxyLogic) BindYxy(in *pb.BindYxyRequest) (*pb.BindYxyResponse, error) {
	deviceID, yxyUID, err := l.svcCtx.Credential.PrepareYxyBind(in.GetDeviceId(), in.GetYxyUid())
	if err != nil {
		if codeErr, ok := errorsx.As(err); ok {
			return nil, codeErr
		}
		l.Errorf("构建用户易校园绑定信息失败: %v", err)
		return nil, errorsx.ErrUnknown
	}

	if err := ensureBindUserExists(l.ctx, l.svcCtx, l.Logger, in.GetUserId()); err != nil {
		return nil, err
	}
	if err := l.svcCtx.UserRepo.UpdateYxyBindByID(l.ctx, in.GetUserId(), deviceID, yxyUID); err != nil {
		l.Errorf("更新用户易校园绑定信息失败: %v", err)
		return nil, errorsx.ErrUnknown
	}

	return &pb.BindYxyResponse{}, nil
}
