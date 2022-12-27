package userService

import (
	"time"
	"usercenter/app/model"
	"usercenter/app/utility"
	"usercenter/config/database"
)

func CreateUser(password, email, sid string) error {
	pass := utility.Encryrpt(password)
	user := &model.User{
		Password:   pass,
		StudentId:  sid,
		Email:      email,
		CreateTime: time.Now(),
		Activate:   0,
	}
	result := database.DB.Create(user)
	if result.Error != nil {
		return result.Error
	} else {
		return nil
	}
}

func CreateUserWithoutEmail(password, email, sid string) error {
	pass := utility.Encryrpt(password)
	user := &model.User{
		Password:   pass,
		StudentId:  sid,
		Email:      email,
		CreateTime: time.Now(),
		Activate:   1,
	}
	result := database.DB.Create(user)
	if result.Error != nil {
		return result.Error
	} else {
		return nil
	}
}
