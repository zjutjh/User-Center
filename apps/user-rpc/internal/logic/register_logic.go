package logic

import (
	"context"
	"errors"
	"strings"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zjutjh/User-Center/apps/user-rpc/internal/dao/model"
	"github.com/zjutjh/User-Center/apps/user-rpc/internal/domain/account"
	"github.com/zjutjh/User-Center/apps/user-rpc/internal/svc"
	"github.com/zjutjh/User-Center/apps/user-rpc/pb"
	"github.com/zjutjh/User-Center/common/errorsx"
)

type RegisterLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterLogic {
	return &RegisterLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *RegisterLogic) Register(in *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	studentID := account.NormalizeStudentID(in.StudentId)
	password := strings.TrimSpace(in.Password)
	cardID := account.NormalizeCardID(in.CardId)
	email := strings.TrimSpace(in.Email)

	if err := account.ValidatePasswordLength(password); err != nil {
		return nil, err
	}

	exists, err := l.svcCtx.UserRepo.ExistsByStudentID(l.ctx, studentID)
	if err != nil {
		l.Errorf("检查用户是否已存在失败: %v", err)
		return nil, errorsx.ErrUnknown
	}
	if exists {
		return nil, errorsx.ErrUserExisted
	}

	student, err := l.svcCtx.StudentRepo.GetStudentByStudentID(l.ctx, studentID)
	if err != nil {
		if codeErr, ok := errorsx.As(err); ok {
			return nil, codeErr
		}
		l.Errorf("查询学生信息失败: %v", err)
		return nil, errorsx.ErrUnknown
	}
	if err := account.VerifyStudentIdentity(student, cardID); err != nil {
		if errors.Is(err, errorsx.ErrParameterInvalid) {
			return nil, errorsx.ErrUserNotExist
		}
		return nil, err
	}

	if err := l.svcCtx.UserRepo.CreateUser(l.ctx, &model.User{
		StudentID: studentID,
		Password:  account.HashPassword(password),
		Email:     email,
	}); err != nil {
		l.Errorf("创建用户失败: %v", err)
		return nil, errorsx.ErrUnknown
	}

	return &pb.RegisterResponse{}, nil
}
