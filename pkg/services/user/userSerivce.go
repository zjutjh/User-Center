package userService

import (
	"errors"
	"log/slog"
	"time"

	"gorm.io/gorm"

	"github.com/zjutjh/User-Center-grpc/pkg/apiExpection"
	"github.com/zjutjh/User-Center-grpc/pkg/database"
	"github.com/zjutjh/User-Center-grpc/pkg/model"
	"github.com/zjutjh/User-Center-grpc/pkg/util"
)

func CheckStudentBySIDAndIID(sid string, iid string) error {
	student := model.Student{}
	result := database.DB.Where(
		&model.Student{
			StudentId: sid,
		},
	).First(&student)
	if student.Iid != iid || result.Error != nil {
		return apiExpection.UserNotFound
	}
	return nil
}

func CreateUser(password, email, sid string) error {
	_, err := GetUserByStudentId(sid)
	if err == nil {
		return apiExpection.UserAlreadyExit
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		slog.Error("failed to get user by student id: %v", err)
		return apiExpection.Unknown
	}
	pass := util.Encryrpt(password)
	user := &model.User{
		Password:   pass,
		StudentId:  sid,
		Email:      email,
		CreateTime: time.Now(),
	}
	return database.DB.Create(user).Error
}

func GetUserByStudentId(studentId string) (*model.User, error) {
	user := model.User{}
	result := database.DB.Where(
		&model.User{
			StudentId: studentId,
		}).First(&user)
	return &user, result.Error
}

func GetUserId(id int) (*model.User, error) {
	user := model.User{}
	result := database.DB.Where(
		&model.User{
			UserId: id,
		},
	).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}

	return &user, nil
}

func UpdateUserPassword(studentId, password string) error {
	user, _ := GetUserByStudentId(studentId)
	pass := util.Encryrpt(password)
	user.Password = pass
	return database.DB.Model(model.User{}).Where(
		model.User{
			StudentId: user.StudentId,
		}).Updates(user).Error
}

func Delete(stuID string) error {
	user, err := GetUserByStudentId(stuID)
	if err != nil {
		return err
	}
	return database.DB.Delete(user).Error
}
