// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package account

import (
	"context"

	"github.com/zjutjh/User-Center/apps/user-api/internal/svc"
	"github.com/zjutjh/User-Center/apps/user-api/internal/types"
	"github.com/zjutjh/User-Center/common/ctxdata"
	"github.com/zjutjh/User-Center/common/errorsx"

	"github.com/zeromicro/go-zero/core/logx"
)

type DelLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 注销账号
func NewDelLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DelLogic {
	return &DelLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DelLogic) Del(req *types.DeleteAccountReq) (resp *types.EmptyResp, err error) {
	userID, ok := ctxdata.UserID(l.ctx)
	if !ok || userID == 0 {
		return nil, errorsx.ErrNotLoggedIn
	}

	if err := l.svcCtx.User.DeleteAccount(l.ctx, userID, req.StudentId, req.IdCard); err != nil {
		l.Errorf("调用用户 RPC 注销账号接口失败: %v", err)
		return nil, err
	}

	return &types.EmptyResp{}, nil
}
