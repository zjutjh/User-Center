package dao

import (
	"context"

	"github.com/zjutjh/User-Center/dao/model"
	"github.com/zjutjh/User-Center/dao/query"
	"github.com/zjutjh/User-Center/pkg/expection"
)

func GetStudentByStudentID(ctx context.Context, studentID string) (*model.Student, error) {
	q := query.Student
	student, err := q.WithContext(ctx).Where(q.StudentID.Eq(studentID)).First()
	if err != nil {
		return nil, expection.UserNotFound
	}
	return student, nil
}
