package logic

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zjutjh/User-Center/apps/user-rpc/internal/domain/account"
	"github.com/zjutjh/User-Center/apps/user-rpc/internal/svc"
	"github.com/zjutjh/User-Center/apps/user-rpc/pb"
	"github.com/zjutjh/User-Center/common/errorsx"
)

type ResetPasswordLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewResetPasswordLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ResetPasswordLogic {
	return &ResetPasswordLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ResetPasswordLogic) ResetPassword(in *pb.ResetPasswordRequest) (*pb.ResetPasswordResponse, error) {
	studentID := account.NormalizeStudentID(in.StudentId)
	cardID := account.NormalizeCardID(in.CardId)
	password := in.Password

	if err := account.ValidatePasswordLength(password); err != nil {
		return nil, err
	}

	user, err := l.svcCtx.UserRepo.GetUserById(l.ctx, in.UserId)
	if err != nil {
		if codeErr, ok := errorsx.As(err); ok {
			return nil, codeErr
		}
		l.Errorf("查询待重置密码用户失败: %v", err)
		return nil, errorsx.ErrUnknown
	}
	if user.StudentID != studentID {
		return nil, errorsx.ErrParameterInvalid
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
		return nil, err
	}

	if err := l.svcCtx.UserRepo.UpdatePasswordByID(l.ctx, in.UserId, account.HashPassword(password)); err != nil {
		l.Errorf("重置用户密码失败: %v", err)
		return nil, errorsx.ErrUnknown
	}

	return &pb.ResetPasswordResponse{}, nil
}
