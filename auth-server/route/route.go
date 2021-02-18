package route

import (
	"auth-server/handler/oauth2_handler"
	"auth-server/handler/user_handler"

	"github.com/kataras/iris/v12"
)

func InitRouter(app *iris.Application) {
	app.PartyFunc("/oauth2", func (u iris.Party) {
		u.Get("/authorize", oauth2_handler.GetAuthorizeHandler)
		u.Post("/authorize", oauth2_handler.PostAuthorizeHandler)
		u.Get("/auth", oauth2_handler.GetAuthHandler)
		u.Post("/auth", oauth2_handler.PostAuthHandler)
		u.Post("/token", oauth2_handler.PostTokenHandler)
		u.Post("/verify", oauth2_handler.PostVerifyTokenHandler)
	})
	app.PartyFunc("/user", func (u iris.Party) {
		u.Post("/register", user_handler.PostRegisterHandler)
		u.Get("/login", user_handler.GetLoginHandler)
		u.Post("/login", user_handler.PostLoginHandler)
		u.Get("/logout", user_handler.GetLogoutHandler)
		u.Post("/password", user_handler.PostSetPasswordHandler)
	})
}