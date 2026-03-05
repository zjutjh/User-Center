package test_test

import (
	"context"
	"flag"
	"log"
	"net"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/zrpc"
	userrpc "github.com/zjutjh/User-Center/apps/user-rpc/app"
	rpcmodel "github.com/zjutjh/User-Center/apps/user-rpc/internal/model"
	"github.com/zjutjh/User-Center/apps/user-rpc/usercenterservice"
)

type rootConfig struct {
	UserRPC userrpc.Config
}

type userCenterRPCTestEnv struct {
	Config    userrpc.Config
	Client    usercenterservice.UserCenterService
	rpcClient zrpc.Client
	server    *zrpc.RpcServer
}

var userCenterRPCTest *userCenterRPCTestEnv

var configFile = flag.String("f", "../../../config.yaml", "配置文件路径")

func TestMain(m *testing.M) {
	os.Exit(UserCenterRPCTestMain(m))
}

func UserCenterRPCTestMain(m *testing.M) int {
	flag.Parse()

	var c rootConfig
	if err := conf.Load(*configFile, &c); err != nil {
		log.Printf("加载配置失败: %v", err)
		return 1
	}

	if err := pingRealMySQL(c.UserRPC); err != nil {
		log.Printf("Ping MySQL 失败: %v", err)
		return 1
	}
	log.Printf("MySQL 连接成功: %s:%d/%s", c.UserRPC.Mysql.Host, c.UserRPC.Mysql.Port, c.UserRPC.Mysql.Database)

	server := userrpc.NewServer(c.UserRPC)
	go server.Start()

	target := clientTarget(c.UserRPC.ListenOn)
	log.Printf("测试进程正在启动 RPC 服务: %s", target)

	rpcClient := zrpc.MustNewClient(zrpc.RpcClientConf{
		Target:   target,
		NonBlock: false,
		Timeout:  3000,
	})

	userCenterRPCTest = &userCenterRPCTestEnv{
		Config:    c.UserRPC,
		Client:    usercenterservice.NewUserCenterService(rpcClient),
		rpcClient: rpcClient,
		server:    server,
	}
	if err := userCenterRPCTest.waitReady(); err != nil {
		log.Printf("等待 RPC 服务就绪失败: %v", err)
		userCenterRPCTest.cleanup()
		return 1
	}
	log.Printf("RPC 服务启动成功: %s", target)

	code := m.Run()
	userCenterRPCTest.cleanup()
	return code
}

func (e *userCenterRPCTestEnv) cleanup() {
	if e.rpcClient != nil {
		_ = e.rpcClient.Conn().Close()
	}
	if e.server != nil {
		e.server.Stop()
	}
}

func (e *userCenterRPCTestEnv) waitReady() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	for {
		_, err := e.Client.HealthyCheck(ctx, &usercenterservice.HealthyCheckRequest{})
		if err == nil {
			return nil
		}

		select {
		case <-ctx.Done():
			return err
		case <-ticker.C:
		}
	}
}

func TestUserCenterServiceHealthyCheck(t *testing.T) {
	ctx := rpcTestContext(t)

	healthyResp, err := userCenterRPCTest.Client.HealthyCheck(ctx, &usercenterservice.HealthyCheckRequest{
		StudentId: "20230001",
	})
	require.NoError(t, err)
	t.Logf("健康检查响应: %+v，错误: %v", healthyResp, err)
}

func TestUserCenterServiceLogin(t *testing.T) {
	ctx := rpcTestContext(t)

	loginResp, err := userCenterRPCTest.Client.Login(ctx, &usercenterservice.LoginRequest{
		StudentId: "302024114514",
		Password:  "secret123",
	})
	t.Logf("登录响应: %+v，错误: %v", loginResp, err)
}

func TestUserCenterServiceRegister(t *testing.T) {
	ctx := rpcTestContext(t)

	registerResp, err := userCenterRPCTest.Client.Register(ctx, &usercenterservice.RegisterRequest{
		StudentId: "302024114514",
		Password:  "114514",
		CardId:    "1145141919166666666",
		Email:     "mjj@zjutjh.com",
	})
	t.Logf("注册响应: %+v，错误: %v", registerResp, err)
}

func TestUserCenterServiceGetUserPassword(t *testing.T) {
	ctx := rpcTestContext(t)

	passwordResp, err := userCenterRPCTest.Client.GetUserPassword(ctx, &usercenterservice.GetUserPasswordRequest{
		UserId: 1,
	})
	t.Logf("获取用户密码响应: %+v，错误: %v", passwordResp, err)
}

func rpcTestContext(t *testing.T) context.Context {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	t.Cleanup(cancel)
	return ctx
}

func pingRealMySQL(c userrpc.Config) error {
	db, err := rpcmodel.NewDB(c.Mysql)
	if err != nil {
		return err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return err
	}
	defer sqlDB.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	return sqlDB.PingContext(ctx)
}

func clientTarget(listenOn string) string {
	host, port, err := net.SplitHostPort(listenOn)
	if err != nil {
		return listenOn
	}
	if host == "" || host == "0.0.0.0" || host == "::" {
		host = "127.0.0.1"
	}
	return net.JoinHostPort(host, port)
}
