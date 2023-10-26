package userService

import (
	"context"
	"time"
	"usercenter/app/apiExpection"
	"usercenter/app/model"
	"usercenter/app/utility"
	"usercenter/config/database"
	"usercenter/config/redis"
)

type UserInRedis struct {
	StudentId string
	Password  string
	Email     string
	Code      string
}

var (
	ctx = context.Background()
)

func CreateUser(password, email, sid string) error {
	pass := utility.Encryrpt(password)
	user := &model.User{
		Password:   pass,
		StudentId:  sid,
		Email:      email,
		CreateTime: time.Now(),
	}
	result := database.DB.Create(user)
	if result.Error != nil {
		return result.Error
	} else {
		return nil
	}
}

func CreateUserInRedis(password, email, sid, code string) error {
	user := &UserInRedis{
		Password:  password,
		StudentId: sid,
		Email:     email,
		Code:      code,
	}
	redis.RedisClient.Set(ctx, user.Email, user, time.Minute*10)
	return nil
}

func CreateUserWithCode(email, code string) error {
	user, err := GetCode(email, code)
	if err != nil {
		return err
	}

	CreateUser(user.Password, user.Email, user.StudentId)

	return nil
}

func GetCode(email, code string) (*UserInRedis, error) {
	var user UserInRedis
	if err := redis.RedisClient.Get(ctx, email).Scan(&user); err != nil {
		return nil, apiExpection.EmailNotFound
	}
	if user.Code != code {
		return nil, apiExpection.CodeError
	}
	return &user, nil
}
