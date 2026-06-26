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

type BindYxyLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 绑定易校园信息
func NewBindYxyLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BindYxyLogic {
	return &BindYxyLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *BindYxyLogic) BindYxy(req *types.BindYxyReq) (resp *types.EmptyResp, err error) {
	userID, ok := ctxdata.UserID(l.ctx)
	if !ok || userID == 0 {
		return nil, errorsx.ErrNotLoggedIn
	}

	if err = l.svcCtx.User.BindYxy(l.ctx, userID, req.DeviceId, req.YxyUid); err != nil {
		l.Errorf("调用用户 RPC 绑定易校园接口失败: %v", err)
		return nil, err
	}

	return &types.EmptyResp{}, nil
}
