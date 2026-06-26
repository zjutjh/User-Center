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

// 统一登录
func OauthHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.LoginReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := authlogic.NewOauthLogic(r.Context(), svcCtx)
		resp, err := l.Oauth(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
