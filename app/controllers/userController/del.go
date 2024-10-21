package userController

import (
	"errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"usercenter/app/apiExpection"
	"usercenter/app/services/studentService"
	"usercenter/app/services/userService"
	"usercenter/app/utility"
)

type form struct {
	IDCard      string `json:"iid" binding:"required"`
	StudentID   string `json:"stuid" binding:"required"`
	BoundSystem uint8  `json:"bound_system"` // 0：wjh 1:foru
}

func DelAccount(c *gin.Context) {
	var postForm form
	err := c.ShouldBindJSON(&postForm)
	if err != nil {
		_ = c.AbortWithError(200, apiExpection.ParamError)
		return
	}
	if !studentService.CheckStudentBYSIDAndIID(postForm.StudentID, postForm.IDCard) {
		utility.JsonResponse(400, "身份证与学号不匹配", nil, c)
		return
	}

	err = userService.DelAccount(postForm.StudentID, postForm.BoundSystem)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		utility.JsonResponse(401, "用户不存在", nil, c)
		return
	} else if !errors.Is(err, errors.New("系统未绑定")) {
		utility.JsonResponse(406, err.Error(), nil, c)
		return
	} else if err != nil {
		_ = c.AbortWithError(200, apiExpection.ServerError)
		return
	}
	utility.JsonSuccessResponse(c, nil)

}
