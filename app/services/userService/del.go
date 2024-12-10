package userService

import (
	"encoding/json"
	"errors"
	"usercenter/app/apiExpection"
	"usercenter/app/model"
	"usercenter/config/database"
	"usercenter/config/redis"
)

func DelAccount(stuID string, system uint8) error {
	user, err := GetUserByStudentId(stuID)
	if err != nil {
		return err
	}
	var boundSystems []model.SystemBinding
	if err := json.Unmarshal(user.BoundSystems, &boundSystems); err != nil {
		return err
	}
	// 删除绑定的系统
	found := false
	for i, binding := range boundSystems {
		if binding.SystemName == model.SystemNameEnum(system) {
			boundSystems = append(boundSystems[:i], boundSystems[i+1:]...)
			found = true
			break
		}
	}
	if !found {
		return errors.New("系统未绑定")
	}
	if len(boundSystems) == 0 {
		result := database.DB.Delete(user)
		return result.Error
	}
	user.BoundSystems, err = json.Marshal(boundSystems)
	if err != nil {
		return err
	}
	result := database.DB.Save(user)
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
