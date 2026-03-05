package main

import (
	"flag"
	"fmt"
	"sync"

	"github.com/zeromicro/go-zero/core/conf"
	userapi "github.com/zjutjh/User-Center/apps/user-api/app"
	userrpc "github.com/zjutjh/User-Center/apps/user-rpc/app"
)

type rootConfig struct {
	UserAPI userapi.Config
	UserRPC userrpc.Config
}

var configFile = flag.String("f", "config.yaml", "the config file")

func main() {
	flag.Parse()

	var c rootConfig
	conf.MustLoad(*configFile, &c)

	rpcServer := userrpc.NewServer(c.UserRPC)
	defer rpcServer.Stop()

	fmt.Printf("用户 RPC 服务启动中，监听地址 %s...\n", c.UserRPC.ListenOn)

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		rpcServer.Start()
	}()

	httpServer := userapi.NewServer(c.UserAPI)
	defer httpServer.Stop()

	fmt.Printf("用户 API 服务启动中，监听地址 %s:%d...\n", c.UserAPI.Host, c.UserAPI.Port)

	go func() {
		defer wg.Done()
		httpServer.Start()
	}()

	wg.Wait()
}
