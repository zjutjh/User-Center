package register

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"time"

	grpcmiddleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpcrecovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/spf13/cobra"
	"github.com/zjutjh/mygo/config"
	"github.com/zjutjh/mygo/nlog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	serverapi "github.com/zjutjh/User-Center/api/user/v1alpha1"
	"github.com/zjutjh/User-Center/biz/handler"
)

// ServerCommandRegister 注册 gRPC + Gateway 服务器启动命令
func ServerCommandRegister() func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		return RunServer()
	}
}

// RunServer 启动 gRPC Server 和 HTTP Gateway
func RunServer() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	grpcAddr := config.Pick().GetString("grpc_server.addr")
	if grpcAddr == "" {
		grpcAddr = ":8080"
	}
	httpAddr := config.Pick().GetString("http_server.addr")
	if httpAddr == "" {
		httpAddr = ":8081"
	}

	// 创建 gRPC Server
	grpcServer := grpc.NewServer(
		grpc.StreamInterceptor(grpcmiddleware.ChainStreamServer(
			grpcrecovery.StreamServerInterceptor(),
		)),
		grpc.UnaryInterceptor(grpcmiddleware.ChainUnaryServer(
			grpcrecovery.UnaryServerInterceptor(),
		)),
	)

	// 注册 gRPC 服务
	serverapi.RegisterUserCenterServiceServer(grpcServer, handler.NewUserHandler())

	// 启动 gRPC Server
	grpcListener, err := net.Listen("tcp", grpcAddr)
	if err != nil {
		return fmt.Errorf("failed to listen gRPC on %s: %w", grpcAddr, err)
	}
	go func() {
		nlog.Pick().Infof("gRPC server listening on %s", grpcAddr)
		if err := grpcServer.Serve(grpcListener); err != nil {
			nlog.Pick().Errorf("gRPC server error: %v", err)
		}
	}()

	// 创建 gRPC-Gateway，通过网络连接独立的 gRPC 端口
	marshaler := &runtime.JSONPb{}
	marshaler.UseProtoNames = false
	marshaler.EmitUnpopulated = true
	gwMux := runtime.NewServeMux(
		runtime.WithMarshalerOption(runtime.MIMEWildcard, marshaler),
	)
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}
	if err := serverapi.RegisterUserCenterServiceHandlerFromEndpoint(ctx, gwMux, grpcAddr, opts); err != nil {
		return fmt.Errorf("failed to register gateway: %w", err)
	}

	// 启动 HTTP Gateway
	nlog.Pick().Infof("HTTP gateway listening on %s", httpAddr)
	httpServer := &http.Server{
		Addr:              httpAddr,
		Handler:           gwMux,
		ReadHeaderTimeout: 60 * time.Second,
	}
	if err := httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("HTTP server error: %w", err)
	}

	return nil
}
