package userController

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"usercenter/app/services/userService"
	"usercenter/app/utility"
	"usercenter/app/utility/oauth"
)

type LoginData struct {
	StudentId   string `json:"stu_id"`
	Password    string `json:"password"`
	BoundSystem uint8  `json:"bound_system"` // 0：wjh 1:foru
}

func AuthPassword(c *gin.Context) {
	var data LoginData
	err := c.ShouldBindJSON(&data)
	if err != nil {
		log.Println(err)
		utility.JsonResponseInternalServerError(c)
		return
	}
	_, err = userService.GetUserByStudentId(data.StudentId)
	if err != nil {
		utility.JsonResponse(404, "该用户不存在", nil, c)
		return
	}
	flag := userService.CheckUserBYStudentIdAndPassword(data.StudentId, data.Password)
	if !flag {
		utility.JsonResponse(409, "密码错误", nil, c)
		return
	}
	err = userService.UpdateBoundSystem(data.StudentId, data.BoundSystem)
	if err != nil {
		utility.JsonResponseInternalServerError(c)
		return
	}
	utility.JsonResponse(200, "OK", nil, c)
}

func OauthPassword(c *gin.Context) {
	var data LoginData
	err := c.ShouldBindJSON(&data)
	if err != nil {
		utility.JsonResponseInternalServerError(c)
		return
	}
	if sid, err := oauth.CheckByOauth(data.StudentId, data.Password); sid != data.StudentId || err != nil {
		fmt.Println(sid, data.StudentId, err)
		if err != nil && err.Error() == "Wrong Password" {
			utility.JsonResponse(409, "密码错误", nil, c)
			return
		} else if err != nil && err.Error() == "Get \"http://www.me.zjut.edu.cn/api/basic/info\": context deadline exceeded (Client.Timeout exceeded while awaiting headers)" {
			utility.JsonResponse(408, "请求超时", nil, c)
			return
		} else {
			utility.JsonResponse(410, "系统异常", nil, c)
			return
		}
	}
	utility.JsonResponse(200, "OK", nil, c)
}
