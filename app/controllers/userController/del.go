package userController

import (
	"errors"
	"usercenter/app/apiExpection"
	"usercenter/app/services/studentService"
	"usercenter/app/services/userService"
	"usercenter/app/utility"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type form struct {
	IDCard    string `json:"iid" binding:"required"`
	StudentID string `json:"stuid" binding:"required"`
}

func DelAccount(c *gin.Context) {
	var postForm form
	err := c.ShouldBindJSON(&postForm)
	if err != nil {
		_ = c.AbortWithError(200, apiExpection.ParamError)
		return
	}
	if !studentService.CheckStudentBYSIDAndIIDLast6(postForm.StudentID, postForm.IDCard) {
		utility.JsonResponse(400, "身份证与学号不匹配", nil, c)
		return
	}

	err = userService.DelAccount(postForm.StudentID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		utility.JsonResponse(401, "用户不存在", nil, c)
		return
	} else if err != nil {
		_ = c.AbortWithError(200, apiExpection.ServerError)
		return
	}
	utility.JsonSuccessResponse(c, nil)

}
