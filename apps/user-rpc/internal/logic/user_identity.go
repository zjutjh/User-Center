package logic

import (
	"context"
	"errors"
	"strings"

	"github.com/zjutjh/User-Center/apps/user-rpc/internal/dao/model"
	"github.com/zjutjh/User-Center/apps/user-rpc/internal/svc"
	"github.com/zjutjh/User-Center/common/errorsx"
	"gorm.io/gorm"
)

func normalizeStudentID(studentID string) string {
	return strings.ToUpper(strings.TrimSpace(studentID))
}

func normalizeCardID(cardID string) string {
	return strings.ToUpper(strings.TrimSpace(cardID))
}

func getUserByID(ctx context.Context, svcCtx *svc.ServiceContext, userID int64) (*model.User, error) {
	user, err := svcCtx.Query.User.WithContext(ctx).
		Where(svcCtx.Query.User.ID.Eq(userID)).
		First()
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errorsx.ErrUserNotExist
	}
	if err != nil {
		return nil, err
	}
	return user, nil
}

func verifyStudentIdentity(ctx context.Context, svcCtx *svc.ServiceContext, studentID, cardID string) error {
	student, err := svcCtx.Query.Student.WithContext(ctx).
		Where(svcCtx.Query.Student.StudentID.Eq(studentID)).
		First()
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return errorsx.ErrUserNotExist
	}
	if err != nil {
		return err
	}
	if !strings.EqualFold(student.IDCard, cardID) {
		return errorsx.ErrParameterInvalid
	}
	return nil
}
