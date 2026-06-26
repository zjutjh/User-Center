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

type BindZfLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 绑定正方密码
func NewBindZfLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BindZfLogic {
	return &BindZfLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *BindZfLogic) BindZf(req *types.BindZfReq) (resp *types.EmptyResp, err error) {
	userID, ok := ctxdata.UserID(l.ctx)
	if !ok || userID == 0 {
		return nil, errorsx.ErrNotLoggedIn
	}

	if err = l.svcCtx.User.BindZf(l.ctx, userID, req.ZfPassword); err != nil {
		l.Errorf("调用用户 RPC 绑定正方接口失败: %v", err)
		return nil, err
	}

	return &types.EmptyResp{}, nil
}
