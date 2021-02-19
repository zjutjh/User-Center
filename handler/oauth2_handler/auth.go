package oauth2_handler

import (
  "github.com/kataras/iris/v12"
  "net/url"
  "strings"
  "ucenter/db"
  "ucenter/model"
  "ucenter/util/http_helper"
)

func GetAuthHandler(ctx iris.Context) {
  sess := http_helper.GetSession(ctx)
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
  sess := http_helper.GetSession(ctx)
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
    ctx.StatusCode(iris.StatusNotAcceptable)
    _, _ = ctx.JSON(model.NewResult(nil, 1101, "missing authorize arguments"))
    return
  }
  ctx.Request().Form = v.(url.Values)
  clientId := ctx.FormValue("client_id")
  userId := uid.(string)
  scopes := strings.Split(ctx.FormValue("scope"), ",")
  db.NewScopeRecord(scopes, clientId, userId)
  PostAuthorizeHandler(ctx)
}
