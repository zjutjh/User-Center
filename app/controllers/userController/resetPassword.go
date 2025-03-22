package userController

import (
	"github.com/gin-gonic/gin"
	"log"
	"usercenter/app/services/studentService"
	"usercenter/app/services/userService"
	"usercenter/app/utility"
)

type ResetPwdForm struct {
	IDCard    string `json:"iid"`
	StudentId string `json:"stuid"`
	Password  string `json:"pwd"`
}

// RePass 不使用密码重置
func RePass(c *gin.Context) {
	var data ResetPwdForm
	err := c.ShouldBindJSON(&data)
	if err != nil {
		log.Println(err)
		utility.JsonResponseInternalServerError(c)
		return
	}

	if len(data.Password) < 6 || len(data.Password) > 20 {
		utility.JsonResponse(401, "密码长度必须在6~20位之间", nil, c)
		return
	}

	if !studentService.CheckStudentBYSIDAndIID(data.StudentId, data.IDCard) {
		utility.JsonResponse(400, "身份证与学号不匹配", nil, c)
		return
	}

	if _, err = userService.GetUserByStudentId(data.StudentId); err != nil {
		utility.JsonResponse(404, "用户不存在", nil, c)
		return
	}

	if err = userService.UpdateUserPasswordByStudentId(data.StudentId, data.Password); err != nil {
		utility.JsonResponseInternalServerError(c)
		return
	}

	utility.JsonResponse(200, "OK", nil, c)
}
