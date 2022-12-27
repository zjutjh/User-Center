package userController

import (
	"github.com/gin-gonic/gin"
	"log"
	"usercenter/app/services/userService"
	"usercenter/app/utility"
)

type LoginData struct {
	StudentId string `json:"stu_id"`
	Password  string `json:"password"`
}

func AuthPassword(c *gin.Context) {
	var data LoginData
	err := c.ShouldBindJSON(&data)
	if err != nil {
		log.Println(err)
		utility.JsonResponseInternalServerError(c)
		return
	}
	user, err1 := userService.GetUserByStudentId(data.StudentId)
	if err1 != nil || user.Activate == 0 {
		utility.JsonResponse(404, "该用户不存在", nil, c)
		return
	}
	flag := userService.CheckUserBYStudentIdAndPassword(data.StudentId, data.Password)
	if !flag {
		utility.JsonResponse(405, "密码错误", nil, c)
		return
	}
	utility.JsonResponse(200, "OK", nil, c)
}
