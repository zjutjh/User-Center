// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package user

import (
	"context"
	"strings"

	"github.com/zjutjh/User-Center/apps/user-api/internal/svc"
	"github.com/zjutjh/User-Center/apps/user-api/internal/types"
	"github.com/zjutjh/User-Center/apps/user-rpc/usercenterservice"
	"github.com/zjutjh/User-Center/common/errorsx"

	"github.com/zeromicro/go-zero/core/logx"
)

type LoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

type LoginResult struct {
	User    *types.UserResp
	Session string
}

// 密码登录
func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LoginLogic) Login(req *types.LoginReq) (resp *LoginResult, err error) {
	studentID := strings.ToUpper(strings.TrimSpace(req.Username))
	password := strings.TrimSpace(req.Password)
	if studentID == "" || password == "" {
		return nil, errorsx.ErrParameterInvalid
	}

	loginResp, err := l.svcCtx.UserRpc.Login(l.ctx, &usercenterservice.LoginRequest{
		StudentId: studentID,
		Password:  password,
	})
	if err != nil {
		l.Errorf("调用用户 RPC 登录接口失败: %v", err)
		return nil, errorsx.FromGRPC(err)
	}
	if loginResp.Session == "" {
		l.Errorf("用户 RPC 登录接口未返回 session, userId=%d", loginResp.UserId)
		return nil, errorsx.ErrUnknown
	}
	sessionUserID, err := l.svcCtx.Session.Decode(loginResp.Session)
	if err != nil || sessionUserID != loginResp.UserId {
		l.Errorf("用户 RPC 登录 session 校验失败: userId=%d sessionUserId=%d err=%v", loginResp.UserId, sessionUserID, err)
		return nil, errorsx.ErrUnknown
	}

	userResp, err := getUserResp(l.ctx, l.svcCtx, loginResp.UserId)
	if err != nil {
		return nil, err
	}

	return &LoginResult{
		User:    userResp,
		Session: loginResp.Session,
	}, nil
}
