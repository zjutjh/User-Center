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
		api.POST("/activation", userController.ActivateUser)
		api.POST("/verify/email", emailController.EmailVerification)
		api.POST("/email", emailController.EmailReset)
		api.POST("/activation/notVerify", userController.ActiviteWithoutEmail)

		//需要邮箱验证
		api.POST("/changePwd", userController.ResetPassword)

		//不需要邮箱验证
		api.POST("/repass", userController.RePass)

		api.POST("/del", userController.DelAccount)
	}
}
