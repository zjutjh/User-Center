package userController

import (
	"github.com/gin-gonic/gin"
	"log"
	"usercenter/app/services/studentService"
	"usercenter/app/services/userService"
	"usercenter/app/utility"
)

func ActiviteWithoutEmail(c *gin.Context) {
	var data RegisterData
	err := c.ShouldBindJSON(&data)
	if err != nil {
		log.Println(err)
		utility.JsonResponseInternalServerError(c)
		return
	}
	if len(data.StudentId) != 12 {
		utility.JsonResponse(402, "学号格式不正确，请重新输入", nil, c)
		return
	}
	_, err = userService.GetUserByStudentId(data.StudentId)
	if err == nil {
		utility.JsonResponse(403, "该通行证已经存在，请重新输入", nil, c)
		return
	}
	flag := studentService.CheckStudentBYSIDAndIID(data.StudentId, data.Iid)
	if !flag {
		utility.JsonResponse(400, "该学号和身份证不存在或者不匹配，请重新输入", nil, c)
		return
	}
	if len(data.Password) < 6 || len(data.Password) > 20 {
		utility.JsonResponse(401, "密码长度必须在6~20位之间", nil, c)
		return
	}
	err = userService.CreateUserWithoutEmail(data.Password, data.Email, data.StudentId)
	if err != nil {
		log.Println(err)
		utility.JsonResponseInternalServerError(c)
		return
	}
	utility.JsonResponse(200, "OK", nil, c)
}
