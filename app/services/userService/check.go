package userService

import (
	"usercenter/app/model"
	"usercenter/app/utility"
	"usercenter/config/database"
)

func CheckUserBYStudentIdAndPassword(studentId, password string) bool {
	pass := utility.Encryrpt(password)
	println(pass)
	user := model.User{}
	result := database.DB.Where(
		model.User{
			StudentId: studentId,
			Password:  pass,
		}).First(&user)
	return result.Error == nil
}
