package userService

import (
	"usercenter/app/model"
	"usercenter/app/utility"
	"usercenter/config/database"
)

func UpdateUserEmailByStudentId(studentId, email string) error {
	user, _ := GetUserByStudentId(studentId)
	user.Email = email
	err := database.DB.Model(model.User{}).Where(
		model.User{
			StudentId: user.StudentId,
		}).Updates(user).Error
	if err != nil {
		return err
	} else {
		return nil
	}
}

func UpdateUserPasswordByStudentId(studentId, password string) error {
	user, _ := GetUserByStudentId(studentId)
	pass := utility.Encryrpt(password)
	user.Password = pass
	err := database.DB.Model(model.User{}).Where(
		model.User{
			StudentId: user.StudentId,
		}).Updates(user).Error
	if err != nil {
		return err
	} else {
		return nil
	}
}
