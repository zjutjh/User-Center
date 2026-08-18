package studentService

import (
	"usercenter/app/model"
	"usercenter/config/database"
)
const iidSuffixLen = 6

func CheckStudentBYSIDAndIIDLast6(sid string, iid string) bool {
	student := model.Student{}
	result := database.DB.Where(
		&model.Student{
			StudentId: sid,
		},
	).First(&student)
	if result.Error != nil || len(iid) < iidSuffixLen || len(student.Iid) < iidSuffixLen {
		return false
	}
	studentIid6 := student.Iid[len(student.Iid)-iidSuffixLen:]
	iid6 := iid[len(iid)-iidSuffixLen:]
	return studentIid6 == iid6
}
