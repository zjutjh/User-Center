package main

import (
  "github.com/go-session/session"
  "github.com/kataras/iris/v12"
  "log"
  "ucenter/config"
  "ucenter/route"
)

func init() {
  session.InitManager(
    session.SetCookieName(config.GetConfig().Session.Name),
    session.SetSign([]byte(config.GetConfig().Session.SecretKey)),
  )
}

func main() {
  app := newApp()
  route.InitRouter(app)
  err := app.Run(iris.Addr(":" + config.GetConfig().Port))
  if err != nil {
    log.Fatalln(err)
  }
}

func newApp() *iris.Application {
  app := iris.New()
  app.RegisterView(iris.HTML("./static", ".html"))
  app.Configure(iris.WithOptimizations)
  app.AllowMethods(iris.MethodOptions)
  return app
}
