package userController

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"log"
	"math/rand"
	"time"
	"usercenter/app/services/emailService"
	"usercenter/app/services/studentService"
	"usercenter/app/services/userService"
	"usercenter/app/utility"
)

type RegisterData struct {
	StudentId   string `json:"stu_id"`
	Password    string `json:"password"`
	Iid         string `json:"iid"`
	Email       string `json:"email"`
	Type        uint8  `json:"type"`         // 0: 本科生 1: 研究生
	BoundSystem uint8  `json:"bound_system"` // 0：wjh 1:foru
}

func ActivateUser(c *gin.Context) {
	var data RegisterData
	err := c.ShouldBindJSON(&data)
	if err != nil {
		log.Println(err)
		utility.JsonResponseInternalServerError(c)
		return
	}

	if !utility.IsValidEmail(data.Email) {
		utility.JsonResponse(405, "邮箱格式不正确", nil, c)
		return
	}

	if len(data.StudentId) != 12 {
		utility.JsonResponse(402, "学号格式不正确，请重新输入", nil, c)
		return
	}
	_, err = userService.GetUserByStudentIdAndSystem(data.StudentId, data.BoundSystem)
	log.Println(err)
	if err == nil {
		utility.JsonResponse(403, "该通行证已经存在，请重新输入", nil, c)
		return
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		_ = c.AbortWithError(200, err)
		return
	}
	flag := studentService.CheckStudentBYSIDAndIID(data.StudentId, data.Iid)
	if !flag {
		utility.JsonResponse(400, "该学号和身份证不存在或者不匹配，请重新输入", nil, c)
		return
	}
	if len(data.Password) < 6 || len(data.Password) > 20 {
		utility.JsonResponse(401, "密码长度必须在6~20位之间", nil, c)
		return
	}

	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	vcode := fmt.Sprintf("%06v", rnd.Int31n(1000000))

	err = userService.CreateUserInRedis(data.Password, data.Email, data.StudentId, vcode, data.Type, data.BoundSystem)
	if err != nil {
		_ = c.AbortWithError(200, err)
		return
	}

	emailService.SendEmail(data.Email, vcode)

	utility.JsonResponse(200, "OK", nil, c)
}
