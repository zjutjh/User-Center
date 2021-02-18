package main

import (
	"auth-server/config"
	"auth-server/route"

	"log"

	"github.com/go-session/session"
	"github.com/kataras/iris/v12"
)

func init() {
	session.InitManager(
		session.SetCookieName(config.Get().Session.Name),
		session.SetSign([]byte(config.Get().Session.SecretKey)),
	)
}

func main() {
	app := newApp()
	route.InitRouter(app)
	err := app.Run(iris.Addr(":" + config.Get().Port))
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