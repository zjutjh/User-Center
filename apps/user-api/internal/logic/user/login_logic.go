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

type LoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 用户登录
func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LoginLogic) Login(req *types.LoginReq) (resp *types.LoginResp, err error) {
	if strings.TrimSpace(req.Username) == "" || strings.TrimSpace(req.Password) == "" {
		return nil, errorsx.ErrParameterInvalid
	}

	rpcResp, err := l.svcCtx.UserRpc.Login(l.ctx, &usercenterservice.LoginRequest{
		StudentId: req.Username,
		Password:  req.Password,
	})
	if err != nil {
		l.Errorf("调用用户 RPC 登录接口失败: %v", err)
		return nil, errorsx.FromGRPC(err)
	}

	return &types.LoginResp{
		UserId: rpcResp.UserId,
	}, nil
}
