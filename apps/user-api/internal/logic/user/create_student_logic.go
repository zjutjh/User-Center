// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package user

import (
	"context"
	"strings"

	"github.com/zjutjh/User-Center/apps/user-api/internal/svc"
	"github.com/zjutjh/User-Center/apps/user-api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateStudentLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 创建学生账号
func NewCreateStudentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateStudentLogic {
	return &CreateStudentLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateStudentLogic) CreateStudent(req *types.CreateStudentReq) (resp *types.CreateStudentResp, err error) {
	studentID := strings.ToUpper(strings.TrimSpace(req.StudentId))
	password := strings.TrimSpace(req.Password)
	cardID := strings.ToUpper(strings.TrimSpace(req.CardId))
	email := strings.TrimSpace(req.Email)

	if err = l.svcCtx.User.Register(l.ctx, studentID, password, cardID, email); err != nil {
		l.Errorf("调用用户 RPC 注册接口失败: %v", err)
		return nil, err
	}

	userID, _, err := l.svcCtx.User.Login(l.ctx, studentID, password)
	if err != nil {
		l.Errorf("注册成功后自动登录失败: %v", err)
		return nil, err
	}

	return &types.CreateStudentResp{
		UserId: userID,
	}, nil
}
