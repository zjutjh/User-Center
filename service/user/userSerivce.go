package userService

import (
	"context"
	"errors"
	"github.com/zjutjh/User-Center-grpc/dao"
	"log/slog"
	"time"

	"gorm.io/gorm"

	"github.com/zjutjh/User-Center-grpc/dao/model"
	"github.com/zjutjh/User-Center-grpc/pkg/expection"
	"github.com/zjutjh/User-Center-grpc/pkg/util"
)

func CheckStudentBySIDAndIID(ctx context.Context, studentID string, iid string) error {
	student, err := dao.GetStudentByStudentID(ctx, studentID)
	if err != nil || student.IDCard != iid {
		return expection.UserNotFound
	}
	return nil
}

func CreateUser(ctx context.Context, password, email, sid string) error {
	user, err := dao.GetUserByStudentId(ctx, sid)
	if user != nil {
		return expection.UserAlreadyExit
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		slog.Error("failed to get user by student id: %v", err)
		return expection.Unknown
	}
	pass := util.Encryrpt(password)
	user = &model.User{
		Password:   pass,
		StudentID:  sid,
		Email:      email,
		CreateTime: time.Now(),
	}
	return dao.CreateUser(ctx, user)
}

func UpdateUserPassword(ctx context.Context, studentId, password string) error {
	user, err := dao.GetUserByStudentId(ctx, studentId)
	if err != nil {
		return err
	}
	pass := util.Encryrpt(password)
	user.Password = pass
	return dao.UpdateUserPassword(ctx, user)
}

func Delete(ctx context.Context, stuID string) error {
	user, err := dao.GetUserByStudentId(ctx, stuID)
	if err != nil {
		return err
	}
	return dao.DeleteUserByStudentId(ctx, user.StudentID)
}

func Login(ctx context.Context, studentID, password string) (*model.User, error) {
	user, err := dao.GetUserByStudentId(ctx, studentID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, expection.UserNotExit
		}
		return nil, expection.Unknown
	}
	if user.Password != util.Encryrpt(password) {
		return nil, expection.AuthError
	}
	return user, nil
}
