package dao

import (
	"context"

	"github.com/zjutjh/User-Center/biz/dao/model"
	"github.com/zjutjh/User-Center/biz/dao/query"
	"github.com/zjutjh/User-Center/biz/err"
)

func GetStudentByStudentID(ctx context.Context, studentID string) (*model.Student, error) {
	q := query.Student
	student, err := q.WithContext(ctx).Where(q.StudentID.Eq(studentID)).First()
	if err != nil {
		return nil, bizerr.UserNotExist
	}
	return student, nil
}
