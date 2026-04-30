// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package user

import (
	"context"

	"github.com/zjutjh/User-Center/apps/user-api/internal/svc"
	"github.com/zjutjh/User-Center/apps/user-api/internal/types"
	"github.com/zjutjh/User-Center/apps/user-rpc/usercenterservice"
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
	userID, err := currentUserID(l.ctx)
	if err != nil {
		return nil, err
	}

	_, err = l.svcCtx.UserRpc.ResetPassword(l.ctx, &usercenterservice.ResetPasswordRequest{
		UserId:    userID,
		StudentId: req.StudentId,
		CardId:    req.IdCard,
		Password:  req.Password,
	})
	if err != nil {
		l.Errorf("调用用户 RPC 重置密码接口失败: %v", err)
		return nil, errorsx.FromGRPC(err)
	}

	return &types.EmptyResp{}, nil
}
