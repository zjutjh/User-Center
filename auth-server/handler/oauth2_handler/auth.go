package oauth2_handler

import (
	"auth-server/database"
	"auth-server/util"
	"github.com/kataras/iris/v12"
	"net/url"
	"strings"
)

func GetAuthHandler(ctx iris.Context) {
	sess := util.GetSession(ctx)
	if sess == nil {
		return
	}
	if _, ok := sess.Get("uid"); !ok {
		ctx.Redirect("/user/login")
		return
	}
	_ = ctx.View("auth.html")
}

func PostAuthHandler(ctx iris.Context) {
	sess := util.GetSession(ctx)
	if sess == nil {
		return
	}
	uid, ok := sess.Get("uid")
	if !ok {
		ctx.Redirect("/user/login")
		return
	}
	v, ok := sess.Get("raw_query_data")
	if !ok {
		_, _ = ctx.WriteString("please authorize again")
		return
	}
	ctx.Request().Form = v.(url.Values)
	clientId := ctx.FormValue("client_id")
	userId := uid.(string)
	scopes := strings.Split(ctx.FormValue("scope"), ",")
	database.NewScopeRecord(scopes, clientId, userId)
	PostAuthorizeHandler(ctx)
}