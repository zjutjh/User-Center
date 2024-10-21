package emailController

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"math/rand"
	"net/http"
	"time"
	"usercenter/app/services/emailService"
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

	if !utility.IsValidEmail(data.Email) {
		utility.JsonResponse(405, "邮箱格式不正确", nil, c)
		return
	}

	err = userService.UpdateUserEmailByStudentId(data.StudentId, data.Email)
	if err != nil {
		log.Println(err)
		utility.JsonResponseInternalServerError(c)
		return
	}

	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	vcode := fmt.Sprintf("%06v", rnd.Int31n(1000000))

	emailService.SendEmail(data.Email, vcode)
	err = userService.CreateUserInRedis("", data.Email, "", vcode, 0, 0)
	if err != nil {
		log.Println(err)
		utility.JsonResponseInternalServerError(c)
		return
	}
	utility.JsonResponse(http.StatusOK, "OK", nil, c)
}
