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
		api.POST("/activation", userController.ActivateUser)
		api.POST("/verify/email", emailController.EmailVerification)
		api.POST("/email", emailController.EmailReset)
		api.POST("/activation/notVerify", userController.ActiviteWithoutEmail)
		api.POST("/changePwd", userController.ResetPassword)
	}
}
