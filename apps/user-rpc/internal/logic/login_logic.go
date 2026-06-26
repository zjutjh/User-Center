package logic

import (
	"context"
	"strings"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zjutjh/User-Center/apps/user-rpc/internal/domain/account"
	"github.com/zjutjh/User-Center/apps/user-rpc/internal/svc"
	"github.com/zjutjh/User-Center/apps/user-rpc/pb"
	"github.com/zjutjh/User-Center/common/errorsx"
)

type LoginLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *LoginLogic) Login(in *pb.LoginRequest) (*pb.LoginResponse, error) {
	studentID := account.NormalizeStudentID(in.StudentId)
	password := strings.TrimSpace(in.Password)

	user, err := l.svcCtx.UserRepo.GetUserByStudentId(l.ctx, studentID)
	if err != nil {
		if codeErr, ok := errorsx.As(err); ok {
			return nil, codeErr
		}
		l.Errorf("根据学号查询用户失败: %v", err)
		return nil, errorsx.ErrUnknown
	}
	if !account.VerifyPassword(user.Password, password) {
		return nil, errorsx.ErrWrongAccountOrPassword
	}

	session, err := l.svcCtx.Session.Encode(user.ID)
	if err != nil {
		l.Errorf("生成登录 session 失败: %v", err)
		return nil, errorsx.ErrUnknown
	}

	return &pb.LoginResponse{
		UserId:  user.ID,
		Session: session,
	}, nil
}
