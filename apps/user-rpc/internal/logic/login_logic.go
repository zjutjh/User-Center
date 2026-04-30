package logic

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"strings"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zjutjh/User-Center/apps/user-rpc/internal/svc"
	"github.com/zjutjh/User-Center/apps/user-rpc/pb"
	"github.com/zjutjh/User-Center/common/errorsx"
	"gorm.io/gorm"
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
	studentID := normalizeStudentID(in.StudentId)
	password := strings.TrimSpace(in.Password)

	if studentID == "" || password == "" {
		return nil, errorsx.ErrWrongAccountOrPassword
	}

	user, err := l.svcCtx.Query.User.WithContext(l.ctx).
		Where(l.svcCtx.Query.User.StudentID.Eq(studentID)).
		First()
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errorsx.ErrUserNotExist
	}
	if err != nil {
		l.Errorf("根据学号查询用户失败: %v", err)
		return nil, errorsx.ErrUnknown
	}
	if user.Password != hashPassword(password) {
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

func hashPassword(password string) string {
	sum := sha256.Sum256([]byte(password))
	return hex.EncodeToString(sum[:])
}
