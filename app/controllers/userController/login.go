package userController

import (
	"github.com/gin-gonic/gin"
	"github.com/zjutjh/WeJH-SDK/oauth"
	"log"
	"usercenter/app/services/userService"
	"usercenter/app/utility"
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
	_, err = oauth.Login(data.StudentId, data.Password)
	if err != nil {
		if err.Error() == "密码错误" {
			utility.JsonResponse(409, "统一系统密码错误", nil, c)
			return
		}
		if err.Error() == "统一系统在夜间关闭" {
			utility.JsonResponse(411, "统一系统在夜间关闭", nil, c)
			return
		}
		if err.Error() == "账号未激活" {
			utility.JsonResponse(412, "统一系统账号未激活", nil, c)
			return
		}
		if err.Error() == "账号错误" {
			utility.JsonResponse(413, "统一系统账号错误", nil, c)
			return
		}
		if err.Error() == "其他错误" {
			utility.JsonResponse(499, "其他错误", nil, c)
			return
		}
	}
	utility.JsonResponse(200, "OK", nil, c)
}
