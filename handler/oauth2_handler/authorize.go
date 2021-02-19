package oauth2_handler

import (
  "github.com/kataras/iris/v12"
  "log"
  "net/url"
  "ucenter/model"
  "ucenter/util/http_helper"
)

func GetAuthorizeHandler(ctx iris.Context) {
  sess := http_helper.GetSession(ctx)
  if sess == nil {
    return
  }
  m := ctx.URLParams()
  clientID := m["client_id"]
  redirectUri := m["redirect_uri"]
  responseType := m["response_type"]
  scope := m["scope"]
  if clientID == "" || redirectUri == "" || responseType == "" || scope == "" {
    ctx.StatusCode(iris.StatusNotAcceptable)
    _, _ = ctx.JSON(model.NewResult(nil, 1101, "missing authorize arguments"))
    return
  }
  if _, err := clientStore.GetByID(nil, clientID); err != nil {
    ctx.StatusCode(iris.StatusNotAcceptable)
    _, _ = ctx.JSON(model.NewResult(nil, 1102, "invalid client_id"))
    return
  }
  if responseType != "code" {
    ctx.StatusCode(iris.StatusNotAcceptable)
    _, _ = ctx.JSON(model.NewResult(nil, 1103, "unsupported response type, only accept authorization_code mode"))
    return
  }
  err := srv.HandleAuthorizeRequest(ctx.ResponseWriter(), ctx.Request())
  if err != nil && err.Error() == "auth_required" {
    return
  }
  if err != nil {
    ctx.StatusCode(iris.StatusNotAcceptable)
    _, _ = ctx.JSON(model.NewResult(nil, 1104, "invalid request: " + err.Error()))
    return
  }
}

func PostAuthorizeHandler(ctx iris.Context) {
  sess := http_helper.GetSession(ctx)
  if sess == nil {
    return
  }
  var form url.Values
  if v, ok := sess.Get("raw_query_data"); ok {
    form = v.(url.Values)
  } else {
    ctx.StatusCode(iris.StatusNotAcceptable)
    _, _ = ctx.JSON(model.NewResult(nil, 1101, "missing authorize arguments"))
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
    ctx.StatusCode(iris.StatusNotAcceptable)
    _, _ = ctx.JSON(model.NewResult(nil, 1104, "invalid request: " + err.Error()))
    return
  }
}
