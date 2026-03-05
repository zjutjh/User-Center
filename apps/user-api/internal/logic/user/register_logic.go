package user

import (
	"context"
	"strings"

	"github.com/zjutjh/User-Center/apps/user-api/internal/svc"
	"github.com/zjutjh/User-Center/apps/user-api/internal/types"
	"github.com/zjutjh/User-Center/apps/user-rpc/usercenterservice"
	"github.com/zjutjh/User-Center/common/errorsx"

	"github.com/zeromicro/go-zero/core/logx"
)

type RegisterLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 用户注册
func NewRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterLogic {
	return &RegisterLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RegisterLogic) Register(req *types.RegisterReq) (resp *types.RegisterResp, err error) {
	if strings.TrimSpace(req.StudentId) == "" || strings.TrimSpace(req.Password) == "" || strings.TrimSpace(req.CardId) == "" {
		return nil, errorsx.ErrParameterInvalid
	}

	_, err = l.svcCtx.UserRpc.Register(l.ctx, &usercenterservice.RegisterRequest{
		StudentId: req.StudentId,
		Password:  req.Password,
		CardId:    req.CardId,
		Email:     req.Email,
	})
	if err != nil {
		l.Errorf("调用用户 RPC 注册接口失败: %v", err)
		return nil, errorsx.FromGRPC(err)
	}

	return &types.RegisterResp{}, nil
}
