package repo

import (
	"context"
	"errors"

	"github.com/zjutjh/User-Center/apps/user-rpc/internal/dao/model"
	"github.com/zjutjh/User-Center/apps/user-rpc/internal/dao/query"
	"github.com/zjutjh/User-Center/common/errorsx"
	"gorm.io/gorm"
)

type StudentRepo struct {
	query *query.Query
}

func NewStudentRepo(q *query.Query) *StudentRepo {
	return &StudentRepo{query: q}
}

func (r *StudentRepo) GetStudentByStudentID(ctx context.Context, studentID string) (*model.Student, error) {
	q := r.query.Student
	student, err := q.WithContext(ctx).Where(q.StudentID.Eq(studentID)).First()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errorsx.ErrUserNotExist
		}
		return nil, err
	}
	return student, nil
}

func (r *StudentRepo) CheckStudentByStudentIDAndCardId(ctx context.Context, studentID string, cardId string) (*model.Student, error) {
	q := r.query.Student
	student, err := q.WithContext(ctx).Where(q.StudentID.Eq(studentID)).First()
	if err != nil || student == nil || student.IDCard != cardId {
		return nil, errorsx.ErrUserNotExist
	}
	return student, nil
}
