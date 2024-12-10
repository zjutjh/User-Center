package userService

import (
	"encoding/json"
	"gorm.io/gorm"
	"usercenter/app/model"
	"usercenter/config/database"
)

func GetUserByEmail(email string) (*model.User, error) {
	user := model.User{}
	result := database.DB.Where(
		&model.User{
			Email: email,
		}).First(&user)
	if result.Error != nil {
		return nil, result.Error
	} else {
		return &user, nil
	}
}

func GetUserByStudentIdAndSystem(studentId string, sid uint8) (*model.User, error) {
	user := model.User{}
	result := database.DB.Where(
		&model.User{
			StudentId: studentId,
		}).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}

	var boundSystems []model.SystemBinding
	if err := json.Unmarshal(user.BoundSystems, &boundSystems); err != nil {
		return nil, err
	}
	for _, system := range boundSystems {
		if system.SystemName == model.SystemNameEnum(sid) {
			return &user, nil
		}
	}
	return nil, gorm.ErrRecordNotFound

}

func GetUserByStudentId(studentId string) (*model.User, error) {
	user := model.User{}
	result := database.DB.Where(
		&model.User{
			StudentId: studentId,
		}).First(&user)
	return &user, result.Error
}

func GetUserId(id int) (*model.User, error) {
	user := model.User{}
	result := database.DB.Where(
		&model.User{
			UserId: id,
		},
	).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}

	return &user, nil
}
