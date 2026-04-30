package logic

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zjutjh/User-Center/apps/user-rpc/internal/dao/model"
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
	studentID := normalizeStudentID(in.StudentId)
	cardID := normalizeCardID(in.CardId)

	if in.UserId <= 0 || studentID == "" || cardID == "" {
		return nil, errorsx.ErrParameterInvalid
	}

	user, err := getUserByID(l.ctx, l.svcCtx, in.UserId)
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

	if err := verifyStudentIdentity(l.ctx, l.svcCtx, studentID, cardID); err != nil {
		if codeErr, ok := errorsx.As(err); ok {
			return nil, codeErr
		}
		l.Errorf("校验注销用户身份失败: %v", err)
		return nil, errorsx.ErrUnknown
	}

	if err := l.svcCtx.DB.WithContext(l.ctx).
		Where("id = ?", in.UserId).
		Delete(&model.User{}).Error; err != nil {
		l.Errorf("删除用户失败: %v", err)
		return nil, errorsx.ErrUnknown
	}

	return &pb.DeleteAccountResponse{}, nil
}
