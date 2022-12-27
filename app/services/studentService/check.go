package studentService

import (
	"usercenter/app/model"
	"usercenter/config/database"
)

func CheckStudentBYSIDAndIID(sid string, iid string) error {
	student := model.Student{}
	result := database.DB.Where(
		&model.Student{
			StudentId: sid,
			Iid:       iid,
		},
	).First(&student)
	println(result.Error)
	if result.Error != nil {
		return result.Error
	} else {
		return nil
	}
}
