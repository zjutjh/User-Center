package user_handler

import (
	"auth-server/database"
	"auth-server/handler/oauth2_handler"
	"auth-server/util"
	"github.com/kataras/iris/v12"
	"log"
)

func GetLoginHandler(ctx iris.Context) {
	sess := util.GetSession(ctx)
	if sess == nil {
		return
	}
	_ = ctx.View("login.html")
}

func PostLoginHandler(ctx iris.Context) {
	sess := util.GetSession(ctx)
	if sess == nil {
		return
	}
	username := ctx.FormValue("username")
	password := ctx.FormValue("password")
	if !util.ValidateUsername(username) || !util.ValidatePassword(password) {
		ctx.StatusCode(iris.StatusNotAcceptable)
		_, _ = ctx.WriteString("invalid username or password")
		return
	}
	ok, res := database.VerifyPassword(username, password)
	if !ok {
		ctx.StatusCode(iris.StatusNotAcceptable)
		_, _ = ctx.WriteString(res)
		return
	}
	sess.Set("uid", ctx.FormValue("username"))
	if _, ok := sess.Get("passive_login"); ok {
		oauth2_handler.PostAuthorizeHandler(ctx)
		sess.Delete("passive_login")
	}
	if err := sess.Save(); err != nil {
		log.Println(err)
	}
}

func GetLogoutHandler(ctx iris.Context) {
	sess := util.GetSession(ctx)
	if sess == nil {
		return
	}
	_ = sess.Flush()
}