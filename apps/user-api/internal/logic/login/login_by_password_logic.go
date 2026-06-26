// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package login

import (
	"context"

	userlogic "github.com/zjutjh/User-Center/apps/user-api/internal/logic/user"
	"github.com/zjutjh/User-Center/apps/user-api/internal/svc"
	"github.com/zjutjh/User-Center/apps/user-api/internal/types"
	"github.com/zjutjh/User-Center/common/errorsx"

	"github.com/zeromicro/go-zero/core/logx"
)

type LoginByPasswordLogic struct {
	logx.Logger
	ctx     context.Context
	svcCtx  *svc.ServiceContext
	session string
}

// 账号密码登录
func NewLoginByPasswordLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginByPasswordLogic {
	return &LoginByPasswordLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LoginByPasswordLogic) Session() string {
	return l.session
}

func (l *LoginByPasswordLogic) LoginByPassword(req *types.LoginReq) (resp *types.UserResp, err error) {
	loginResp, err := userlogic.NewLoginLogic(l.ctx, l.svcCtx).Login(req)
	if err != nil {
		return nil, err
	}
	if loginResp.Session == "" {
		l.Errorf("账号密码登录未生成 session")
		return nil, errorsx.ErrUnknown
	}

	l.session = loginResp.Session
	return loginResp.User, nil
}
