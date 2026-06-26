// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package bind

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"github.com/zjutjh/User-Center/apps/user-api/internal/logic/bind"
	"github.com/zjutjh/User-Center/apps/user-api/internal/svc"
	"github.com/zjutjh/User-Center/apps/user-api/internal/types"
)

// 绑定统一认证密码
func BindOauthHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.BindOauthReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := bind.NewBindOauthLogic(r.Context(), svcCtx)
		resp, err := l.BindOauth(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
