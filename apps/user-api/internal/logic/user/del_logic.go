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
	userID, err := currentUserID(l.ctx)
	if err != nil {
		return nil, err
	}

	_, err = l.svcCtx.UserRpc.DeleteAccount(l.ctx, &usercenterservice.DeleteAccountRequest{
		UserId:    userID,
		StudentId: req.StudentId,
		CardId:    req.IdCard,
	})
	if err != nil {
		l.Errorf("调用用户 RPC 注销账号接口失败: %v", err)
		return nil, errorsx.FromGRPC(err)
	}

	return &types.EmptyResp{}, nil
}
