package svc

import (
	"github.com/zeromicro/go-zero/rest"
	"github.com/zjutjh/User-Center/apps/user-api/internal/config"
	"github.com/zjutjh/User-Center/apps/user-api/internal/domain"
	"github.com/zjutjh/User-Center/apps/user-api/internal/middleware"
	"github.com/zjutjh/User-Center/common/sessionx"
)

type ServiceContext struct {
	Config   config.Config
	UserAuth rest.Middleware
	User     *domain.UserCenter
	Session  *sessionx.Manager
}

func NewServiceContext(c config.Config) *ServiceContext {
	session := sessionx.NewManager(c.Session)

	return &ServiceContext{
		Config:   c,
		UserAuth: middleware.NewUserAuthMiddleware(session).Handle,
		User:     domain.NewUserCenter(c.UserRpc),
		Session:  session,
	}
}
