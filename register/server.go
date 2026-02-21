package register

import (
	"net"

	//grpcmiddleware "github.com/grpc-ecosystem/go-grpc-middleware"
	//grpcrecovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"github.com/spf13/cobra"
	"github.com/zjutjh/User-Center/api/rpc"
	"github.com/zjutjh/mygo/config"
	"github.com/zjutjh/mygo/nlog"
	"google.golang.org/grpc"

	serverapi "github.com/zjutjh/User-Center/idl/user/v1alpha1"
)

// GRPCServerCommandRegister 注册 gRPC 服务器启动命令
func GRPCServerCommandRegister() func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		return RunGRPCServer()
	}
}

// RunGRPCServer 启动 gRPC Server
func RunGRPCServer() error {
	grpcAddr := config.Pick().GetString("grpc_server.addr")
	if grpcAddr == "" {
		grpcAddr = ":8080"
	}

	// 创建 gRPC Server
	grpcServer := grpc.NewServer(
	//grpc.UnaryInterceptor(grpcmiddleware.ChainUnaryServer(
	//	grpcrecovery.UnaryServerInterceptor(),
	//)),
	)

	// 注册 gRPC 服务
	serverapi.RegisterUserCenterServiceServer(grpcServer, rpc.NewUserService())

	// 启动 gRPC Server
	grpcListener, err := net.Listen("tcp", grpcAddr)
	if err != nil {
		return err
	}
	nlog.Pick().Infof("gRPC 监听端口 %s", grpcAddr)
	if err = grpcServer.Serve(grpcListener); err != nil {
		return err
	}

	return nil
}
