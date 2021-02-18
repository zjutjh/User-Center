package user_handler

import (
	"auth-server/database"
	"auth-server/util"
	"github.com/kataras/iris/v12"
)

// only for test purpose
func PostRegisterHandler(ctx iris.Context) {
	sess := util.GetSession(ctx)
	if sess == nil {
		return
	}
	if ctx.Request().Form != nil {
		if err := ctx.Request().ParseForm(); err != nil {
			ctx.StatusCode(iris.StatusInternalServerError)
			_, _ = ctx.WriteString("invalid form")
			return
		}
	}
	username := ctx.FormValue("username")
	password := ctx.FormValue("password")
	if !util.ValidateUsername(username) || !util.ValidatePassword(password) {
		ctx.StatusCode(iris.StatusNotAcceptable)
		_, _ = ctx.WriteString("invalid username or password")
		return
	}
	database.NewUser(username, password)
	ctx.StatusCode(iris.StatusOK)
}
