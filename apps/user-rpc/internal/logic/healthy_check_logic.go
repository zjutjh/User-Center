package logic

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zjutjh/User-Center/apps/user-rpc/internal/svc"
	"github.com/zjutjh/User-Center/apps/user-rpc/pb"
)

type HealthyCheckLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewHealthyCheckLogic(ctx context.Context, svcCtx *svc.ServiceContext) *HealthyCheckLogic {
	return &HealthyCheckLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *HealthyCheckLogic) HealthyCheck(_ *pb.HealthyCheckRequest) (*pb.HealthyCheckResponse, error) {
	return &pb.HealthyCheckResponse{}, nil
}
