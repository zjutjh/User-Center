package dao

import (
	"context"

	"github.com/zjutjh/User-Center/biz/dao/model"
	"github.com/zjutjh/User-Center/biz/dao/query"
	"github.com/zjutjh/User-Center/biz/err"
)

func GetUserByStudentId(ctx context.Context, studentId string) (*model.User, error) {
	q := query.User
	user, err := q.WithContext(ctx).Where(q.StudentID.Eq(studentId)).First()
	if err != nil {
		return nil, bizerr.UserNotExist
	}
	return user, nil
}

func UpdateUserPassword(ctx context.Context, user *model.User) error {
	q := query.User
	_, err := q.WithContext(ctx).Where(q.ID.Eq(user.ID)).Updates(user)
	return err
}

func DeleteUserByStudentId(ctx context.Context, studentId string) error {
	q := query.User
	_, err := q.WithContext(ctx).Where(q.StudentID.Eq(studentId)).Delete()
	return err
}

func CreateUser(ctx context.Context, user *model.User) error {
	q := query.User
	return q.WithContext(ctx).Create(user)
}
