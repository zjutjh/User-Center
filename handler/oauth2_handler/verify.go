package oauth2_handler

import (
  "github.com/kataras/iris/v12"
  "strings"
  "ucenter/db"
  "ucenter/model"
)

func VerifyTokenHandler(ctx iris.Context) (bool, *model.Result) {
  _, err := srv.ValidationBearerToken(ctx.Request())
  if err != nil {
    return false, model.NewResult(nil, 1106, "invalid access token")
  }
  userId := ctx.FormValue("username")
  clientId := ctx.FormValue("client_id")
  permission := ctx.FormValue("scope")
  if userId == "" {
    return false, model.NewResult(nil, 1107, "invalid username")
  }
  if !db.ValidationTokenScope(strings.Split(permission, ","), clientId, userId) {
    return false, model.NewResult(nil, 1108, "invalid scope")
  }
  return true, nil
}
