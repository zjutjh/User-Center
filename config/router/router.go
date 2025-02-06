package router

import (
	"github.com/gin-gonic/gin"
	"usercenter/app/controllers/emailController"
	"usercenter/app/controllers/userController"
)

func Init(r *gin.Engine) {

	const pre = "/api"

	api := r.Group(pre)
	{
		api.POST("/auth", userController.AuthPassword)
		api.POST("/oauth", userController.OauthPassword)
		api.POST("/email", emailController.EmailReset)
		api.POST("/activation/notVerify", userController.ActiviteWithoutEmail)

		//不需要邮箱验证
		api.POST("/repass", userController.RePass)

		api.POST("/del", userController.DelAccount)
	}
}
