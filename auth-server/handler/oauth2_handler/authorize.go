package oauth2_handler

import (
	"auth-server/util"
	"github.com/kataras/iris/v12"
	"log"
	"net/url"
)

func GetAuthorizeHandler(ctx iris.Context) {
	sess := util.GetSession(ctx)
	if sess == nil {
		return
	}
	m := ctx.URLParams()
	clientID := m["client_id"]
	redirectUri := m["redirect_uri"]
	responseType := m["response_type"]
	scope := m["scope"]
	if clientID == "" || redirectUri == "" || responseType == "" || scope == "" {
		_, _ = ctx.WriteString("missing arguments")
		return
	}
	if _, err := clientStore.GetByID(nil, clientID); err != nil {
		_, _ = ctx.WriteString("invalid client_id")
		return
	}
	if responseType != "code" {
		_, _ = ctx.WriteString("unsupported response type")
		return
	}
	err := srv.HandleAuthorizeRequest(ctx.ResponseWriter(), ctx.Request())
	if err != nil && err.Error() == "auth_required" {
		return
	}
	if err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		_, _ = ctx.WriteString("invalid request")
	}
}

func PostAuthorizeHandler(ctx iris.Context) {
	sess := util.GetSession(ctx)
	if sess == nil {
		return
	}
	var form url.Values
	if v, ok := sess.Get("raw_query_data"); ok {
		form = v.(url.Values)
	} else {
		ctx.StatusCode(iris.StatusBadRequest)
		_, _ = ctx.WriteString("missing arguments")
		return
	}
	ctx.Request().Form = form
	sess.Delete("raw_query_data")
	if err := sess.Save(); err != nil {
		log.Println(err)
	}
	err := srv.HandleAuthorizeRequest(ctx.ResponseWriter(), ctx.Request())
	if err != nil && err.Error() == "auth_required" {
		return
	}
	if err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		_, _ = ctx.WriteString("invalid request")
		return
	}
}