package route

import (
  "github.com/kataras/iris/v12"
  "ucenter/handler/api_handler"
  "ucenter/handler/oauth2_handler"
  "ucenter/handler/user_handler"
)

func InitRouter(app *iris.Application) {
  app.PartyFunc("/oauth2", func (u iris.Party) {
    u.Get("/authorize", oauth2_handler.GetAuthorizeHandler)
    u.Post("/authorize", oauth2_handler.PostAuthorizeHandler)
    u.Get("/auth", oauth2_handler.GetAuthHandler)
    u.Post("/auth", oauth2_handler.PostAuthHandler)
    u.Post("/token", oauth2_handler.PostTokenHandler)
  })
  app.PartyFunc("/user", func (u iris.Party) {
    u.Post("/register", user_handler.PostRegisterHandler)
    u.Get("/login", user_handler.GetLoginHandler)
    u.Post("/login", user_handler.PostLoginHandler)
    u.Get("/logout", user_handler.GetLogoutHandler)
    u.Post("/password", user_handler.PostSetPasswordHandler)
  })
  app.PartyFunc("/api", func (u iris.Party) {
    u.Get("/data", api_handler.GetDataHandler)
  })
}
