package userController

import (
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
		utility.JsonResponse(405, "密码错误", nil, c)
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
		log.Println(err)
		utility.JsonResponseInternalServerError(c)
		return
	}
	if sid, err := oauth.CheckByOauth(data.StudentId, data.Password); sid != data.StudentId || err != nil {
		log.Println(err)
		utility.JsonResponse(405, "密码错误", nil, c)
		return
	}
	utility.JsonResponse(200, "OK", nil, c)
}
