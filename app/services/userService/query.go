package userService

import (
	"usercenter/app/model"
	"usercenter/config/database"
)

func GetUserByEmail(email string) (*model.User, error) {
	user := model.User{}
	result := database.DB.Where(
		&model.User{
			Email: email,
		}).First(&user)
	if result.Error != nil {
		return nil, result.Error
	} else {
		return &user, nil
	}
}

func GetUserByStudentId(studentId string) (*model.User, error) {
	user := model.User{}
	result := database.DB.Where(
		&model.User{
			StudentId: studentId,
		}).First(&user)
	if result.Error != nil {
		return nil, result.Error
	} else {
		return &user, nil
	}
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
