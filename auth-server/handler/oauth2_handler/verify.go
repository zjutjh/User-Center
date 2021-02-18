package oauth2_handler

import (
	"auth-server/database"
	"github.com/kataras/iris/v12"
	"strings"
)

func PostVerifyTokenHandler(ctx iris.Context) {
	_, _, err := srv.ValidationTokenRequest(ctx.Request())
	if err != nil {
		ctx.StatusCode(iris.StatusNotAcceptable)
		_, _ = ctx.WriteString("invalid token")
		return
	}
	userId := ctx.FormValue("username")
	clientId := ctx.FormValue("client_id")
	permission := ctx.FormValue("scope")
	if userId == "" {
		ctx.StatusCode(iris.StatusNotAcceptable)
		_, _ = ctx.WriteString("invalid username")
		return
	}
	if !database.ValidationTokenScope(strings.Split(permission, ","), clientId, userId) {
		ctx.StatusCode(iris.StatusNotAcceptable)
		_, _ = ctx.WriteString("invalid permission scope")
		return
	}
	ctx.StatusCode(iris.StatusOK)
	_, _ = ctx.WriteString(ctx.FormValue("scope"))
	return
}
