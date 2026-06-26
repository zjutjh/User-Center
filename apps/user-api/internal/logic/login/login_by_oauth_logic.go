// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package login

import (
	"context"

	"github.com/zjutjh/User-Center/apps/user-api/internal/svc"
	"github.com/zjutjh/User-Center/apps/user-api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type LoginByOauthLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 统一登录
func NewLoginByOauthLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginByOauthLogic {
	return &LoginByOauthLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LoginByOauthLogic) LoginByOauth(req *types.LoginReq) (resp *types.UserResp, err error) {
	// todo: add your logic here and delete this line

	return
}
