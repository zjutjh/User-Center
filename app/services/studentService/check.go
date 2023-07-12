package studentService

import (
	"usercenter/app/model"
	"usercenter/config/database"
)

func CheckStudentBYSIDAndIID(sid string, iid string) bool {
	student := model.Student{}
	result := database.DB.Where(
		&model.Student{
			StudentId: sid,
		},
	).First(&student)
	if student.Iid != iid || result.Error != nil {
		return false
	} else {
		return true
	}
}
