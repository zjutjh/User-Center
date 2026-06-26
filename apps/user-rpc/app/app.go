package userrpc

import (
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"github.com/zjutjh/User-Center/apps/user-rpc/internal/config"
	"github.com/zjutjh/User-Center/apps/user-rpc/internal/server"
	"github.com/zjutjh/User-Center/apps/user-rpc/internal/svc"
	"github.com/zjutjh/User-Center/apps/user-rpc/pb"
	"github.com/zjutjh/User-Center/common/interceptors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type Config = config.Config

func NewServer(c Config) *zrpc.RpcServer {
	ctx := svc.NewServiceContext(c)

	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		pb.RegisterUserCenterServiceServer(grpcServer, server.NewUserCenterServiceServer(ctx))

		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})

	s.AddUnaryInterceptors(interceptors.UnaryErrorInterceptor())

	return s
}
