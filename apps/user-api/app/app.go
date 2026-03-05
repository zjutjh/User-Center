package userapi

import (
	"github.com/zeromicro/go-zero/rest"
	"github.com/zjutjh/User-Center/apps/user-api/internal/config"
	"github.com/zjutjh/User-Center/apps/user-api/internal/handler"
	"github.com/zjutjh/User-Center/apps/user-api/internal/svc"
	"github.com/zjutjh/User-Center/common/resultx"
)

type Config = config.Config

func NewServer(c Config) *rest.Server {
	resultx.InstallHTTPHandlers()

	server := rest.MustNewServer(c.RestConf)
	ctx := svc.NewServiceContext(c)
	handler.RegisterHandlers(server, ctx)

	return server
}
