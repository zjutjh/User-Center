package svc

import (
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
	"github.com/zjutjh/User-Center/apps/user-api/internal/config"
	"github.com/zjutjh/User-Center/apps/user-api/internal/middleware"
	"github.com/zjutjh/User-Center/apps/user-rpc/usercenterservice"
	"github.com/zjutjh/User-Center/common/sessionx"
)

type ServiceContext struct {
	Config   config.Config
	UserAuth rest.Middleware
	UserRpc  usercenterservice.UserCenterService
	Session  *sessionx.Manager
}

func NewServiceContext(c config.Config) *ServiceContext {
	client := zrpc.MustNewClient(c.UserRpc)
	session := sessionx.NewManager(c.Session)

	return &ServiceContext{
		Config:   c,
		UserAuth: middleware.NewUserAuthMiddleware(session).Handle,
		UserRpc:  usercenterservice.NewUserCenterService(client),
		Session:  session,
	}
}
