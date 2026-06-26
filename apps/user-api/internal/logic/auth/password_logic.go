// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package auth

import (
	"context"
	"strings"

	"github.com/zjutjh/User-Center/apps/user-api/internal/svc"
	"github.com/zjutjh/User-Center/apps/user-api/internal/types"
	"github.com/zjutjh/User-Center/common/errorsx"

	"github.com/zeromicro/go-zero/core/logx"
)

type PasswordLogic struct {
	logx.Logger
	ctx     context.Context
	svcCtx  *svc.ServiceContext
	session string
}

type loginResult struct {
	User    *types.UserResp
	Session string
}

// 账号密码登录
func NewPasswordLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PasswordLogic {
	return &PasswordLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PasswordLogic) Session() string {
	return l.session
}

func (l *PasswordLogic) Password(req *types.LoginReq) (resp *types.UserResp, err error) {
	loginResp, err := l.login(req)
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

func (l *PasswordLogic) login(req *types.LoginReq) (*loginResult, error) {
	studentID := strings.ToUpper(strings.TrimSpace(req.Username))
	password := strings.TrimSpace(req.Password)

	userID, session, err := l.svcCtx.User.Login(l.ctx, studentID, password)
	if err != nil {
		l.Errorf("调用用户 RPC 登录接口失败: %v", err)
		return nil, err
	}
	if session == "" {
		l.Errorf("用户 RPC 登录接口未返回 session, userId=%d", userID)
		return nil, errorsx.ErrUnknown
	}
	sessionUserID, err := l.svcCtx.Session.Decode(session)
	if err != nil || sessionUserID != userID {
		l.Errorf("用户 RPC 登录 session 校验失败: userId=%d sessionUserId=%d err=%v", userID, sessionUserID, err)
		return nil, errorsx.ErrUnknown
	}

	userResp, err := l.svcCtx.User.GetUser(l.ctx, userID)
	if err != nil {
		return nil, err
	}

	return &loginResult{
		User:    userResp,
		Session: session,
	}, nil
}
