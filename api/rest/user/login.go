package user

import (
	"reflect"
	"runtime"

	"github.com/gin-gonic/gin"
	"github.com/zjutjh/User-Center/comm"
	pb "github.com/zjutjh/User-Center/idl/user/v1alpha1"
	usercenter "github.com/zjutjh/User-Center/mygo_plugin"
	"github.com/zjutjh/mygo/foundation/reply"
	"github.com/zjutjh/mygo/kit"
	"github.com/zjutjh/mygo/nlog"
	"github.com/zjutjh/mygo/session"
	"github.com/zjutjh/mygo/swagger"
)

// LoginHandler API router注册点
func LoginHandler() gin.HandlerFunc {
	api := LoginApi{}
	swagger.CM[runtime.FuncForPC(reflect.ValueOf(hfLogin).Pointer()).Name()] = api
	return hfLogin
}

type LoginApi struct {
	Info     struct{}         `name:"用户登录" desc:"通过学号密码登录"`
	Request  LoginApiRequest  // API请求参数 (Uri/Header/Query/Body)
	Response LoginApiResponse // API响应数据 (Body中的Data部分)
}

type LoginApiRequest struct {
	Body struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
}

type LoginApiResponse struct {
	Header struct {
		Session string `json:"session"`
	}
}

// Run Api业务逻辑执行点
func (u *LoginApi) Run(ctx *gin.Context) kit.Code {
	// 调用 gRPC Login 服务
	client := usercenter.Pick()
	resp, err := client.Login(ctx.Request.Context(), &pb.LoginRequest{
		StudentId: u.Request.Body.Username,
		Password:  u.Request.Body.Password,
	})
	if err != nil {
		nlog.Pick().WithContext(ctx).WithError(err).Error("gRPC 调用失败")
		return comm.CodeThirdServiceError
	}
	if resp.Code == pb.BizCode_OK {
		if err := session.SetIdentity(ctx, resp.UserId); err != nil {
			nlog.Pick().WithContext(ctx).WithError(err).Error("设置 session 失败")
			return comm.CodeMiddlewareServiceError
		}
		return comm.CodeOK
	}
	return comm.FromBizCode(resp.Code)
}

// Init Api初始化 进行参数校验和绑定
func (u *LoginApi) Init(ctx *gin.Context) (err error) {
	return ctx.ShouldBindJSON(&u.Request.Body)
}

// hfLogin API执行入口
func hfLogin(ctx *gin.Context) {
	api := &LoginApi{}
	err := api.Init(ctx)
	if err != nil {
		nlog.Pick().WithContext(ctx).WithError(err).Warn("参数绑定校验错误")
		reply.Fail(ctx, comm.CodeParameterInvalid)
		return
	}
	code := api.Run(ctx)
	if !ctx.IsAborted() {
		if code == comm.CodeOK {
			reply.Success(ctx, api.Response)
		} else {
			reply.Fail(ctx, code)
		}
	}
}
