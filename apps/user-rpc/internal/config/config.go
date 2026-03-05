package config

import (
	"github.com/zeromicro/go-zero/zrpc"
	"github.com/zjutjh/User-Center/common/dbx"
)

type Config struct {
	zrpc.RpcServerConf
	Mysql dbx.MysqlConf
}
