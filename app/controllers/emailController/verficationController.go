package emailController

import (
	"github.com/gin-gonic/gin"
	"log"
	"usercenter/app/services/userService"
	"usercenter/app/utility"
)

type VerificationData struct {
	Email            string `json:"email"`
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

	err = userService.CreateUserWithCode(data.Email, data.VerificationCode)
	if err != nil {
		_ = c.AbortWithError(200, err)
		return
	}

	utility.JsonResponse(200, "OK", nil, c)
}
