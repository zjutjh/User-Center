package config

import (
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
	"github.com/zjutjh/User-Center/common/sessionx"
)

type Config struct {
	rest.RestConf
	UserRpc zrpc.RpcClientConf
	Session sessionx.Config
}
