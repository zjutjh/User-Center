package oauth2_handler

import (
  "encoding/json"
  "github.com/kataras/iris/v12"
  "ucenter/model"
  "ucenter/util/http_helper"
)

func PostTokenHandler(ctx iris.Context) {
  grantType := ctx.FormValue("grant_type")
  if grantType != "authorization_code" {
    ctx.StatusCode(iris.StatusNotAcceptable)
    _, _ = ctx.JSON(model.NewResult(nil, 1103,"unsupported grant type, only accept authorization_code mode"))
    return
  }
  rw := http_helper.NewFakeResponse()
  err := srv.HandleTokenRequest(rw, ctx.Request())
  if err != nil {
    ctx.StatusCode(iris.StatusNotAcceptable)
    _, _ = ctx.JSON(model.NewResult(nil, 1104, "invalid request: " + err.Error()))
    return
  }
  var rwJson model.OAuth2TokenResult
  _ = json.Unmarshal(rw.GetBody(), &rwJson)
  if rwJson.Error != "" {
    ctx.StatusCode(iris.StatusNotAcceptable)
    _, _ = ctx.JSON(model.NewResult(rwJson, 1105, "internal error: " + rwJson.Error))
  } else {
    ctx.StatusCode(iris.StatusOK)
    _, _ = ctx.JSON(model.NewResult(rwJson, 0, "successfully fetched token"))
  }
}
