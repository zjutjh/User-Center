package oauth2_handler

import (
	"github.com/kataras/iris/v12"
)

func PostTokenHandler(ctx iris.Context) {
	err := srv.HandleTokenRequest(ctx.ResponseWriter(), ctx.Request())
	if err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		_, _ = ctx.WriteString("internal error: token handler failed")
		return
	}
}