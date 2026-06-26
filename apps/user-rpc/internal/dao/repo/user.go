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
	user, err := q.WithContext(ctx).Where(q.ID.Eq(userId)).First()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errorsx.ErrUserNotExist
		}
		return nil, err
	}
	return user, nil
}

func (r *UserRepo) ExistsByStudentID(ctx context.Context, studentID string) (bool, error) {
	q := r.query.User
	_, err := q.WithContext(ctx).Where(q.StudentID.Eq(studentID)).First()
	if err == nil {
		return true, nil
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false, nil
	}
	return false, err
}

func (r *UserRepo) UpdatePasswordByID(ctx context.Context, userID int64, password string) error {
	q := r.query.User
	_, err := q.WithContext(ctx).Where(q.ID.Eq(userID)).Update(q.Password, password)
	return err
}

func (r *UserRepo) UpdateYxyBindByID(ctx context.Context, userID int64, deviceID, yxyUID string) error {
	q := r.query.User
	_, err := q.WithContext(ctx).Where(q.ID.Eq(userID)).Updates(&model.User{
		DeviceID: deviceID,
		YxyUID:   yxyUID,
	})
	return err
}

func (r *UserRepo) UpdateZfBindByID(ctx context.Context, userID int64, zfPassword string) error {
	q := r.query.User
	_, err := q.WithContext(ctx).Where(q.ID.Eq(userID)).Update(q.ZfPassword, zfPassword)
	return err
}

func (r *UserRepo) UpdateOauthBindByID(ctx context.Context, userID int64, oauthPassword string) error {
	q := r.query.User
	_, err := q.WithContext(ctx).Where(q.ID.Eq(userID)).Update(q.OauthPassword, oauthPassword)
	return err
}

func (r *UserRepo) DeleteUserByID(ctx context.Context, userID int64) error {
	q := r.query.User
	_, err := q.WithContext(ctx).Where(q.ID.Eq(userID)).Delete()
	return err
}
