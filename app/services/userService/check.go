package userService

import (
	"encoding/json"
	"time"
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

func UpdateBoundSystem(studentId string, boundSystem uint8) error {
	user, err := GetUserByStudentId(studentId)
	if err != nil {
		return err
	}
	var boundSystems []model.SystemBinding
	if err := json.Unmarshal(user.BoundSystems, &boundSystems); err != nil {
		return err
	}
	for _, system := range boundSystems {
		if system.SystemName == model.SystemNameEnum(boundSystem) {
			return nil
		}
	}
	boundSystems = append(boundSystems, model.SystemBinding{
		SystemName: model.SystemNameEnum(boundSystem),
		BoundAt:    time.Now(),
	})
	user.BoundSystems, err = json.Marshal(boundSystems)
	if err != nil {
		return err
	}
	result := database.DB.Save(user)
	return result.Error
}
