// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package auth

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	authlogic "github.com/zjutjh/User-Center/apps/user-api/internal/logic/auth"
	"github.com/zjutjh/User-Center/apps/user-api/internal/svc"
	"github.com/zjutjh/User-Center/apps/user-api/internal/types"
)

// 账号密码登录
func PasswordHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.LoginReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := authlogic.NewPasswordLogic(r.Context(), svcCtx)
		resp, err := l.Password(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			svcCtx.Session.SetEncoded(w, l.Session())
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
