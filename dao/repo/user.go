package repo

import (
	"context"
	"errors"

	"github.com/zjutjh/User-Center/comm"
	"github.com/zjutjh/User-Center/dao/model"
	"github.com/zjutjh/User-Center/dao/query"
	"github.com/zjutjh/mygo/ndb"
	"gorm.io/gorm"
)

type UserRepo struct {
	query *query.Query
}

func NewUserRepo(tx ...*query.Query) *UserRepo {
	var q *query.Query
	if len(tx) > 0 {
		q = tx[0]
	} else {
		q = query.Use(ndb.Pick())
	}
	return &UserRepo{
		query: q,
	}
}

func (r *UserRepo) GetUserByStudentId(ctx context.Context, studentId string) (*model.User, error) {
	q := r.query.User
	user, err := q.WithContext(ctx).Where(q.StudentID.Eq(studentId)).First()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, comm.UserNotExist
		} else {
			return nil, err
		}
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
