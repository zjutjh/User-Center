package userController

import (
	"usercenter/app/services/studentService"
	"usercenter/app/services/userService"
	"usercenter/app/utility"

	"github.com/gin-gonic/gin"
)

type RegisterData struct {
	StudentId string `json:"stu_id"`
	Password  string `json:"password"`
	Iid       string `json:"iid"`
	Email     string `json:"email"`
	Type      uint8  `json:"type"` // 0: 本科生 1: 研究生
}

func Activite(c *gin.Context) {
	var data RegisterData
	err := c.ShouldBindJSON(&data)
	if err != nil {
		utility.JsonResponse(400, "参数错误", nil, c)
		return
	}
	if len(data.Password) < 6 || len(data.Password) > 20 {
		utility.JsonResponse(401, "密码长度必须在6~20位之间", nil, c)
		return
	}
	flag := studentService.CheckStudentBYSIDAndIIDLast6(data.StudentId, data.Iid)
	if !flag {
		utility.JsonResponse(400, "该学号和身份证不存在或者不匹配，请重新输入", nil, c)
		return
	}
	err = userService.CreateUser(data.Password, data.Email, data.StudentId, data.Type)
	if err != nil && err.Error() == "密码错误" {
		utility.JsonResponse(407, "密码错误", nil, c)
		return
	} else if err != nil {
		utility.JsonResponseInternalServerError(c)
		return
	}
	utility.JsonResponse(200, "OK", nil, c)
}
