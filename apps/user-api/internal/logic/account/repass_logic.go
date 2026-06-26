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

type RepassLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 重置密码
func NewRepassLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RepassLogic {
	return &RepassLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RepassLogic) Repass(req *types.ResetPasswordReq) (resp *types.EmptyResp, err error) {
	userID, ok := ctxdata.UserID(l.ctx)
	if !ok || userID == 0 {
		return nil, errorsx.ErrNotLoggedIn
	}

	if err := l.svcCtx.User.ResetPassword(l.ctx, userID, req.StudentId, req.IdCard, req.Password); err != nil {
		l.Errorf("调用用户 RPC 重置密码接口失败: %v", err)
		return nil, err
	}

	return &types.EmptyResp{}, nil
}
