package main

import (
	"flag"
	"fmt"

	"github.com/zeromicro/go-zero/core/conf"
	userapi "github.com/zjutjh/User-Center/apps/user-api/app"
)

type rootConfig struct {
	UserAPI userapi.Config
}

var configFile = flag.String("f", "config.yaml", "the config file")

func main() {
	flag.Parse()

	var c rootConfig
	conf.MustLoad(*configFile, &c)
	server := userapi.NewServer(c.UserAPI)
	defer server.Stop()

	fmt.Printf("用户 API 服务启动中，监听地址 %s:%d...\n", c.UserAPI.Host, c.UserAPI.Port)
	server.Start()
}
