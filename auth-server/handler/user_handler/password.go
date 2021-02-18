package user_handler

import (
	"auth-server/database"
	"auth-server/util"
	"github.com/kataras/iris/v12"
)

func PostSetPasswordHandler(ctx iris.Context) {
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
	oldPassword := ctx.FormValue("old_password")
	newPassword := ctx.FormValue("new_password")
	if !util.ValidateUsername(username) || !util.ValidatePassword(oldPassword) || !util.ValidatePassword(newPassword) {
		ctx.StatusCode(iris.StatusNotAcceptable)
		_, _ = ctx.WriteString("invalid username or password")
		return
	}
	ok, res := database.VerifyPassword(username, oldPassword)
	if !ok {
		ctx.StatusCode(iris.StatusNotAcceptable)
		_, _ = ctx.WriteString(res)
		return
	}
	database.SetPassword(username, newPassword)
}
