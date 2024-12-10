package userService

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"gorm.io/gorm"
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
	Type      uint8
	System    uint8
	Code      string
}

var (
	ctx = context.Background()
)

func CreateUser(password, email, sid string, userType uint8, bound uint8) error {
	user, err := GetUserByStudentId(sid)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	pass := utility.Encryrpt(password)
	var boundSystems []model.SystemBinding
	if err == nil {
		flag := CheckUserBYStudentIdAndPassword(sid, password)
		if !flag {
			return errors.New("密码错误")
		}
		if err := json.Unmarshal(user.BoundSystems, &boundSystems); err != nil {
			return err
		}
	} else {
		user = &model.User{
			Password:   pass,
			StudentId:  sid,
			Email:      email,
			Type:       userType,
			CreateTime: time.Now(),
		}
	}
	boundSystems = append(boundSystems, model.SystemBinding{
		SystemName: model.SystemNameEnum(bound),
		BoundAt:    time.Now(),
	})

	user.BoundSystems, err = json.Marshal(boundSystems)
	if err != nil {
		return err
	}

	if user.UserId == 0 {
		result := database.DB.Create(user)
		return result.Error
	} else {
		result := database.DB.Save(user)
		return result.Error
	}
}

func CreateUserInRedis(password, email, sid, code string, userType uint8, bound uint8) error {
	user := &UserInRedis{
		Password:  password,
		StudentId: sid,
		Email:     email,
		Code:      code,
		Type:      userType,
		System:    bound,
	}
	userData, err := json.Marshal(user)
	if err != nil {
		return fmt.Errorf("json解析失败: %v", err)
	}

	statusCmd := redis.RedisClient.Set(ctx, user.Email, userData, time.Minute*10)
	if err := statusCmd.Err(); err != nil {
		return fmt.Errorf("验证码存储redis失败: %v", err)
	}
	return nil
}

func CreateUserWithCode(email, code string) error {
	user, err := GetCode(email, code)
	if err != nil {
		return err
	}

	err = CreateUser(user.Password, user.Email, user.StudentId, user.Type, user.System)
	return err
}

func GetCode(email, code string) (*UserInRedis, error) {
	// 从 Redis 中获取用户数据
	userData, err := redis.RedisClient.Get(ctx, email).Result()
	if err != nil {
		return nil, apiExpection.EmailNotFound
	}
	// 将 JSON 数据反序列化为 UserInRedis 结构体
	var user UserInRedis
	if err := json.Unmarshal([]byte(userData), &user); err != nil {
		return nil, fmt.Errorf("failed to unmarshal user data: %v", err)
	}
	// 检查提供的 code 是否匹配
	if user.Code != code {
		return nil, apiExpection.CodeError
	}
	err = DelCode(email)
	return &user, err
}
