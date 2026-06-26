package logic

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zjutjh/User-Center/apps/user-rpc/internal/domain/account"
	"github.com/zjutjh/User-Center/apps/user-rpc/internal/svc"
	"github.com/zjutjh/User-Center/apps/user-rpc/pb"
	"github.com/zjutjh/User-Center/common/errorsx"
)

type DeleteAccountLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteAccountLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteAccountLogic {
	return &DeleteAccountLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeleteAccountLogic) DeleteAccount(in *pb.DeleteAccountRequest) (*pb.DeleteAccountResponse, error) {
	studentID := account.NormalizeStudentID(in.StudentId)
	cardID := account.NormalizeCardID(in.CardId)

	user, err := l.svcCtx.UserRepo.GetUserById(l.ctx, in.UserId)
	if err != nil {
		if codeErr, ok := errorsx.As(err); ok {
			return nil, codeErr
		}
		l.Errorf("查询待注销用户失败: %v", err)
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

	if err := l.svcCtx.UserRepo.DeleteUserByID(l.ctx, in.UserId); err != nil {
		l.Errorf("删除用户失败: %v", err)
		return nil, errorsx.ErrUnknown
	}

	return &pb.DeleteAccountResponse{}, nil
}
