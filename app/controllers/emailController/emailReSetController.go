package emailController

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"usercenter/app/services/emailService"
	"usercenter/app/services/redisService"
	"usercenter/app/services/userService"
	"usercenter/app/utility"
)

type EmailData struct {
	Email     string `json:"email"`
	StudentId string `json:"stu_id"`
}

func EmailReset(c *gin.Context) {
	var data EmailData
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
	err = userService.UpdateUserEmailByStudentId(data.StudentId, data.Email)
	if err != nil {
		log.Println(err)
		utility.JsonResponseInternalServerError(c)
		return
	}
	code := emailService.SendEmail(data.Email)
	if code != "" {
		redisService.SetRedis(code, data.Email)
	}
	utility.JsonResponse(http.StatusOK, "OK", nil, c)
}
