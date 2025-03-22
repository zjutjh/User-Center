package userController

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/zjutjh/WeJH-SDK/oauth"
	"github.com/zjutjh/WeJH-SDK/oauth/oauthException"
	"log"
	"usercenter/app/apiExpection"
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
	_, e := oauth.Login(data.StudentId, data.Password)
	if e != nil {
		switch {
		case errors.Is(e, oauthException.ClosedError):
			_ = c.AbortWithError(200, apiExpection.ClosedError)
		case errors.Is(e, oauthException.WrongPassword):
			_ = c.AbortWithError(200, apiExpection.WrongPassword)
		case errors.Is(e, oauthException.NotActivatedError):
			_ = c.AbortWithError(200, apiExpection.NotActivatedError)
		case errors.Is(e, oauthException.WrongAccount):
			_ = c.AbortWithError(200, apiExpection.WrongAccount)
		case errors.Is(e, oauthException.OtherError):
			_ = c.AbortWithError(200, apiExpection.OtherError("其他错误"))
		}
		return
	}
	utility.JsonSuccessResponse(c, nil)
}
