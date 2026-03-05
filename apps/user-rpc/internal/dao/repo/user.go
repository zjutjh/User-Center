package repo

import (
	"context"
	"errors"

	"github.com/zjutjh/User-Center/apps/user-rpc/internal/dao/model"
	"github.com/zjutjh/User-Center/apps/user-rpc/internal/dao/query"
	"github.com/zjutjh/User-Center/common/errorsx"
	"gorm.io/gorm"
)

type UserRepo struct {
	query *query.Query
}

func NewUserRepo(q *query.Query) *UserRepo {
	return &UserRepo{query: q}
}

func (r *UserRepo) GetUserByStudentId(ctx context.Context, studentId string) (*model.User, error) {
	q := r.query.User
	user, err := q.WithContext(ctx).Where(q.StudentID.Eq(studentId)).First()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errorsx.ErrUserNotExist
		}
		return nil, err
	}
	return user, nil
}

func (r *UserRepo) UpdateUserPassword(ctx context.Context, user *model.User) error {
	q := r.query.User
	_, err := q.WithContext(ctx).Where(q.ID.Eq(user.ID)).Updates(user)
	return err
}

func (r *UserRepo) DeleteUserByStudentId(ctx context.Context, studentId string) error {
	q := r.query.User
	_, err := q.WithContext(ctx).Where(q.StudentID.Eq(studentId)).Delete()
	return err
}

func (r *UserRepo) CreateUser(ctx context.Context, user *model.User) error {
	q := r.query.User
	return q.WithContext(ctx).Create(user)
}

func (r *UserRepo) GetUserById(ctx context.Context, userId int64) (*model.User, error) {
	q := r.query.User
	return q.WithContext(ctx).Where(q.ID.Eq(userId)).First()
}
