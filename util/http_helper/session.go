package http_helper

import (
	"github.com/go-session/session"
  "github.com/kataras/iris/v12"
  "log"
)

func GetSession(ctx iris.Context) session.Store {
  sess, err := session.Start(ctx.Request().Context(), ctx.ResponseWriter(), ctx.Request())
  if err != nil {
    ctx.StatusCode(iris.StatusInternalServerError)
    if _, err = ctx.WriteString("session internal error"); err != nil {
      log.Println(err)
    }
    return nil
  }
  return sess
}
