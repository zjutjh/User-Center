package userService

import (
	"usercenter/app/apiExpection"
	"usercenter/config/database"
	"usercenter/config/redis"
)

func DelAccount(stuID string) error {
	user, err := GetUserByStudentId(stuID)
	if err != nil {
		return err
	}
	result := database.DB.Delete(user)
	return result.Error
}

func DelCode(email string) error {
	// 从 Redis 中获取用户数据
	result, err := redis.RedisClient.Del(ctx, email).Result()
	if err != nil {
		return err
	}
	if result == 0 {
		return apiExpection.EmailNotFound
	}
	return nil
}
