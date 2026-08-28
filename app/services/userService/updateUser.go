package userService

import (
	"usercenter/app/apiExpection"
	"usercenter/app/model"
	"usercenter/app/utility"
	"usercenter/config/database"
)

func UpdateUserEmailByStudentId(studentId, email string) error {
	user, err := GetUserByStudentId(studentId)
	if err != nil {
		return err
	}
	if user.Email == email {
		return nil
	}
	result := database.DB.Model(&model.User{}).
		Where("student_id = ?", studentId).
		Update("email", email)

	return result.Error
}

func UpdateUserPasswordByStudentId(studentId, password string) error {
	user, err := GetUserByStudentId(studentId)
	if err != nil {
		return err
	}
	pass := utility.Encryrpt(password)
	if user.Password == pass {
		return apiExpection.UpdateSame
	}
	result := database.DB.Model(&model.User{}).
		Where("student_id = ?", studentId).
		Update("password", pass)

	return result.Error
}
