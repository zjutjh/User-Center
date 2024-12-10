package emailController

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
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

	utility.JsonResponse(http.StatusOK, "OK", nil, c)
}
