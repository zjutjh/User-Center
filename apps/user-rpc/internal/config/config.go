package config

import (
	"github.com/zeromicro/go-zero/zrpc"
	"github.com/zjutjh/User-Center/common/dbx"
	"github.com/zjutjh/User-Center/common/sessionx"
)

type Config struct {
	zrpc.RpcServerConf
	Mysql   dbx.MysqlConf
	Session sessionx.Config
}
