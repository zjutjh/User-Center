// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package login

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"github.com/zjutjh/User-Center/apps/user-api/internal/logic/login"
	"github.com/zjutjh/User-Center/apps/user-api/internal/svc"
	"github.com/zjutjh/User-Center/apps/user-api/internal/types"
)

// 账号密码登录
func LoginByPasswordHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.LoginReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := login.NewLoginByPasswordLogic(r.Context(), svcCtx)
		resp, err := l.LoginByPassword(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			svcCtx.Session.SetEncoded(w, l.Session())
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
