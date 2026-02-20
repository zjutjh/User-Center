package main

import (
	"github.com/spf13/cobra"
	"github.com/zjutjh/mygo/foundation/command"

	"github.com/zjutjh/User-Center/register"
)

func main() {
	command.Execute(
		register.Boot,    // 应用引导注册器
		register.Command, // 应用命令注册器
		func(cmd *cobra.Command, args []string) error {
			// 默认启动 gRPC + HTTP Gateway 服务
			return register.RunServer()
		},
	)
}
