package register

import (
	"fmt"
	"net"
	"os"
	"time"

	"github.com/spf13/cobra"
	"github.com/zjutjh/User-Center/api/rpc"
	"github.com/zjutjh/mygo/config"
	"github.com/zjutjh/mygo/foundation/kernel"
	"google.golang.org/grpc"

	serverapi "github.com/zjutjh/User-Center/idl/user/v1alpha1"
)

type Config struct {
	Addr                string        `mapstructure:"addr"`
	ShutdownWaitTimeout time.Duration `mapstructure:"shutdown_wait_timeout"`
}

var DefaultConfig = Config{
	Addr:                ":8080",
	ShutdownWaitTimeout: 10 * time.Second,
}

// GrpcServerCommandRegister 注册 gRPC 服务器启动命令
func GrpcServerCommandRegister() func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		StartGrpcServer()
		return nil
	}
}

func StartGrpcServer() {
	// 获取配置
	conf := DefaultConfig
	_ = config.Pick().UnmarshalKey("grpc_server", &conf)

	grpcServer := grpc.NewServer()
	// 注册路由
	serviceRegister(grpcServer)

	// 启动http server
	go listenAndServeGrpcServer(grpcServer, conf.Addr)

	// 监听等待关闭服务
	kernel.ListenStop(func() error {
		grpcServer.GracefulStop()
		_, _ = fmt.Fprintln(os.Stdout, "Grpc Server关闭完成")
		return nil
	})
}

func serviceRegister(grpcServer *grpc.Server) {
	serverapi.RegisterUserCenterServiceServer(grpcServer, &rpc.UserService{})
}
func listenAndServeGrpcServer(grpcServer *grpc.Server, addr string) {
	grpcListener, err := net.Listen("tcp", addr)
	if err != nil {
		_, _ = fmt.Fprintln(os.Stdout, "监听 gRPC 端口失败:", err)
	}
	if err = grpcServer.Serve(grpcListener); err != nil {
		_, _ = fmt.Fprintln(os.Stdout, "gRPC 服务启动失败:", err)
	}
}

//
//// StartGrpcServer 启动 gRPC Server
//func StartGrpcServer() {
//	grpcAddr := config.Pick().GetString("grpc_server.addr")
//	if grpcAddr == "" {
//		grpcAddr = ":8080"
//	}
//
//	// 创建 gRPC Server
//	grpcServer := grpc.NewServer()
//
//	// 注册 gRPC 服务
//
//	// 启动 gRPC Server
//	grpcListener, err := net.Listen("tcp", grpcAddr)
//	if err != nil {
//		nlog.Pick().WithError(err).Error("监听 gRPC 端口失败")
//		os.Exit(1)
//	}
//	nlog.Pick().Infof("gRPC 监听端口 %s", grpcAddr)
//	if err = grpcServer.Serve(grpcListener); err != nil {
//		nlog.Pick().WithError(err).Error("gRPC 服务启动失败")
//		os.Exit(1)
//	}
//}
