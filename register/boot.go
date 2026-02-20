package register

import (
	"github.com/zjutjh/mygo/foundation/kernel"
	"github.com/zjutjh/mygo/ndb"
	"github.com/zjutjh/mygo/nedis"
	"github.com/zjutjh/mygo/nlog"

	"github.com/zjutjh/User-Center/biz/dao/query"
)

func Boot() kernel.BootList {
	return kernel.BootList{
		// 基础引导器
		nlog.Boot(), // 业务日志

		// Client引导器
		ndb.Boot(),   // DB
		nedis.Boot(), // Redis

		// 业务引导器
		InitQueryBoot(), // 初始化 gen query
	}
}

// InitQueryBoot 初始化 gorm/gen 的全局 query 对象
func InitQueryBoot() func() error {
	return func() error {
		query.SetDefault(ndb.Pick())
		return nil
	}
}
