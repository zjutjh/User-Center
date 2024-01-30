package userService

import "usercenter/config/database"

func DelAccount(stuID string) error {
	user, err := GetUserByStudentId(stuID)
	if err != nil {
		return err
	}
	result := database.DB.Delete(user)
	return result.Error
}
