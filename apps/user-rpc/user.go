package main

import (
	"flag"
	"fmt"

	"github.com/zeromicro/go-zero/core/conf"
	userrpc "github.com/zjutjh/User-Center/apps/user-rpc/app"
)

type rootConfig struct {
	UserRPC userrpc.Config
}

var configFile = flag.String("f", "config.yaml", "the config file")

func main() {
	flag.Parse()

	var c rootConfig
	conf.MustLoad(*configFile, &c)
	s := userrpc.NewServer(c.UserRPC)
	defer s.Stop()

	fmt.Printf("用户 RPC 服务启动中，监听地址 %s...\n", c.UserRPC.ListenOn)
	s.Start()
}
