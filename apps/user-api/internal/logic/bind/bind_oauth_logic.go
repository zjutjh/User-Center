// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package bind

import (
	"context"

	"github.com/zjutjh/User-Center/apps/user-api/internal/svc"
	"github.com/zjutjh/User-Center/apps/user-api/internal/types"
	"github.com/zjutjh/User-Center/common/ctxdata"
	"github.com/zjutjh/User-Center/common/errorsx"

	"github.com/zeromicro/go-zero/core/logx"
)

type BindOauthLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 绑定统一认证密码
func NewBindOauthLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BindOauthLogic {
	return &BindOauthLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *BindOauthLogic) BindOauth(req *types.BindOauthReq) (resp *types.EmptyResp, err error) {
	userID, ok := ctxdata.UserID(l.ctx)
	if !ok || userID == 0 {
		return nil, errorsx.ErrNotLoggedIn
	}

	if err = l.svcCtx.User.BindOauth(l.ctx, userID, req.OauthPassword); err != nil {
		l.Errorf("调用用户 RPC 绑定统一认证接口失败: %v", err)
		return nil, err
	}

	return &types.EmptyResp{}, nil
}
