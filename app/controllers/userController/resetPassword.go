package userController

import (
	"github.com/gin-gonic/gin"
	"log"
	"usercenter/app/services/redisService"
	"usercenter/app/services/userService"
	"usercenter/app/utility"
)

type ResetPwdData struct {
	Password  string `json:"pwd"`
	StudentId string `json:"stu_id"`
	Code      string `json:"code"`
}

func ResetPassword(c *gin.Context) {
	var data ResetPwdData
	err := c.ShouldBindJSON(&data)
	if err != nil {
		log.Println(err)
		utility.JsonResponseInternalServerError(c)
		return
	}
	_, err = userService.GetUserByStudentId(data.StudentId)
	if err != nil {
		log.Println(err)
		utility.JsonResponse(404, "该用户不存在", nil, c)
		return
	}
	email := redisService.GetRedis(data.Code)
	if email == "" {
		utility.JsonResponse(406, "验证码错误", nil, c)
		return
	}
	if len(data.Password) < 6 || len(data.Password) > 20 {
		utility.JsonResponse(401, "密码长度必须在6~20位之间", nil, c)
		return
	}
	err = userService.UpdateUserPasswordByStudentIdAndEmail(data.StudentId, data.Password, email)
	if err != nil {
		log.Println(err)
		utility.JsonResponseInternalServerError(c)
		return
	}
	utility.JsonResponse(200, "OK", nil, c)
}
