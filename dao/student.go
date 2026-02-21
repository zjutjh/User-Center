package dao

import (
	"context"

	"github.com/zjutjh/User-Center/comm/bizerr"
	"github.com/zjutjh/User-Center/dao/model"
	"github.com/zjutjh/User-Center/dao/query"
)

func GetStudentByStudentID(ctx context.Context, studentID string) (*model.Student, error) {
	q := query.Student
	student, err := q.WithContext(ctx).Where(q.StudentID.Eq(studentID)).First()
	if err != nil {
		return nil, bizerr.UserNotExist
	}
	return student, nil
}
