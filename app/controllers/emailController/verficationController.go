package emailController

import (
	"github.com/gin-gonic/gin"
	"log"
	"usercenter/app/services/redisService"
	"usercenter/app/services/userService"
	"usercenter/app/utility"
)

type VerificationData struct {
	VerificationCode string `json:"verification_code"`
}

func EmailVerification(c *gin.Context) {
	var data VerificationData
	err := c.ShouldBindJSON(&data)
	if err != nil {
		log.Println(err)
		utility.JsonResponseInternalServerError(c)
		return
	}
	email := redisService.GetRedis(data.VerificationCode)
	if email == "" {
		utility.JsonResponse(406, "验证码错误", nil, c)
		return
	}
	user, _ := userService.GetUserByEmail(email)
	err = userService.Activate(email, user)
	if err != nil {
		log.Println(err)
		utility.JsonResponseInternalServerError(c)
		return
	}
	utility.JsonResponse(200, "OK", nil, c)
}
