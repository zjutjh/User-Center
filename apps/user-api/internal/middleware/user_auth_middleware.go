package middleware

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"github.com/zjutjh/User-Center/common/ctxdata"
	"github.com/zjutjh/User-Center/common/sessionx"
)

type UserAuthMiddleware struct {
	session *sessionx.Manager
}

func NewUserAuthMiddleware(session *sessionx.Manager) *UserAuthMiddleware {
	return &UserAuthMiddleware{session: session}
}

func (m *UserAuthMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, err := m.session.UserID(r)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		next(w, r.WithContext(ctxdata.WithUserID(r.Context(), userID)))
	}
}
