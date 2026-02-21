package register

import (
	"github.com/spf13/cobra"
	"github.com/zjutjh/mygo/foundation/command"
)

func Command(root *cobra.Command) {
	// 业务命令
	command.Add("server", GRPCServerCommandRegister())
}
