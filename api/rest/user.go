package rest

import (
	"reflect"
	"runtime"

	"github.com/gin-gonic/gin"
	"github.com/zjutjh/User-Center/comm"
	"github.com/zjutjh/mygo/foundation/reply"
	"github.com/zjutjh/mygo/kit"
	"github.com/zjutjh/mygo/nlog"
	"github.com/zjutjh/mygo/swagger"
)

// UserHandler API router注册点
func UserHandler() gin.HandlerFunc {
	api := UserApi{}
	swagger.CM[runtime.FuncForPC(reflect.ValueOf(hfUser).Pointer()).Name()] = api
	return hfUser
}

type UserApi struct {
	Info     struct{}        `name:"API名称" desc:"API描述"`
	Request  UserApiRequest  // API请求参数 (Uri/Header/Query/Body)
	Response UserApiResponse // API响应数据 (Body中的Data部分)
}

type UserApiRequest struct {
}

type UserApiResponse struct{}

// Run Api业务逻辑执行点
func (u *UserApi) Run(ctx *gin.Context) kit.Code {
	// TODO: 在此处编写接口业务逻辑
	return comm.CodeOK
}

// Init Api初始化 进行参数校验和绑定
func (u *UserApi) Init(ctx *gin.Context) (err error) {
	return err
}

// hfUser API执行入口
func hfUser(ctx *gin.Context) {
	api := &UserApi{}
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
