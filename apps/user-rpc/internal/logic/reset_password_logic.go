package logic

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zjutjh/User-Center/apps/user-rpc/internal/dao/model"
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
	studentID := normalizeStudentID(in.StudentId)
	cardID := normalizeCardID(in.CardId)
	password := in.Password

	if in.UserId <= 0 || studentID == "" || cardID == "" || password == "" {
		return nil, errorsx.ErrParameterInvalid
	}
	if len(password) < 6 || len(password) > 20 {
		return nil, errorsx.ErrPasswordLength
	}

	user, err := getUserByID(l.ctx, l.svcCtx, in.UserId)
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

	if err := verifyStudentIdentity(l.ctx, l.svcCtx, studentID, cardID); err != nil {
		if codeErr, ok := errorsx.As(err); ok {
			return nil, codeErr
		}
		l.Errorf("校验学生身份失败: %v", err)
		return nil, errorsx.ErrUnknown
	}

	if err := l.svcCtx.DB.WithContext(l.ctx).
		Model(&model.User{}).
		Where("id = ?", in.UserId).
		Update("password", hashPassword(password)).Error; err != nil {
		l.Errorf("重置用户密码失败: %v", err)
		return nil, errorsx.ErrUnknown
	}

	return &pb.ResetPasswordResponse{}, nil
}
