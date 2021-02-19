package user_handler

import (
  "github.com/kataras/iris/v12"
  "log"
  "ucenter/db"
  "ucenter/handler/oauth2_handler"
  "ucenter/model"
  "ucenter/util/http_helper"
  "ucenter/util/validation"
)

func GetLoginHandler(ctx iris.Context) {
  sess := http_helper.GetSession(ctx)
  if sess == nil {
    return
  }
  _ = ctx.View("login.html")
}

func PostLoginHandler(ctx iris.Context) {
  sess := http_helper.GetSession(ctx)
  if sess == nil {
    return
  }
  username := ctx.FormValue("username")
  password := ctx.FormValue("password")
  if !validation.ValidateUsername(username) || !validation.ValidatePassword(password) {
    ctx.StatusCode(iris.StatusNotAcceptable)
    _, _ = ctx.JSON(model.NewResult(nil, 1001, "invalid username or password"))
    return
  }
  ok, res := db.VerifyPassword(username, password)
  if !ok {
    ctx.StatusCode(iris.StatusNotAcceptable)
    _, _ = ctx.JSON(res)
    return
  }
  sess.Set("uid", ctx.FormValue("username"))
  if _, ok := sess.Get("passive_login"); ok {
    oauth2_handler.PostAuthorizeHandler(ctx)
    sess.Delete("passive_login")
    if err := sess.Save(); err != nil {
      log.Println(err)
    }
    return
  }
  if err := sess.Save(); err != nil {
    log.Println(err)
  }
  ctx.StatusCode(iris.StatusOK)
  _, _ = ctx.JSON(model.NewResult(nil, 0, "successfully logged in"))
}

func GetLogoutHandler(ctx iris.Context) {
  sess := http_helper.GetSession(ctx)
  if sess == nil {
    return
  }
  _ = sess.Flush()
  ctx.Redirect("/")
}
