// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package auth

import (
	"context"

	"github.com/zjutjh/User-Center/apps/user-api/internal/svc"
	"github.com/zjutjh/User-Center/apps/user-api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type OauthLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 统一登录
func NewOauthLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OauthLogic {
	return &OauthLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *OauthLogic) Oauth(req *types.LoginReq) (resp *types.UserResp, err error) {
	// todo: add your logic here and delete this line

	return
}
