package user_handler

import (
  "github.com/kataras/iris/v12"
  "ucenter/db"
  "ucenter/model"
  "ucenter/util/http_helper"
  "ucenter/util/validation"
)

func PostRegisterHandler(ctx iris.Context) {
  sess := http_helper.GetSession(ctx)
  if sess == nil {
    return
  }
  if ctx.Request().Form != nil {
    if err := ctx.Request().ParseForm(); err != nil {
      ctx.StatusCode(iris.StatusNotAcceptable)
      _, _ = ctx.JSON(model.NewResult(nil, 1003, "invalid form submitted"))
      return
    }
  }
  username := ctx.FormValue("username")
  password := ctx.FormValue("password")
  if !validation.ValidateUsername(username) || !validation.ValidatePassword(password) {
    ctx.StatusCode(iris.StatusNotAcceptable)
    _, _ = ctx.JSON(model.NewResult(nil, 1001, "invalid username or password"))
    return
  }
  db.NewUser(username, password)
  ctx.StatusCode(iris.StatusOK)
  _, _ = ctx.JSON(model.NewResult(nil, 0, "successfully registered"))
}
