package userService

import (
	"usercenter/app/model"
	"usercenter/config/database"
)

func Activate(email string, user *model.User) error {
	user.Activate = 1
	err := database.DB.Model(model.User{}).Where(
		model.User{
			StudentId: user.StudentId,
			Email:     email,
		}).Updates(user).Error
	if err != nil {
		return err
	} else {
		return nil
	}
}
