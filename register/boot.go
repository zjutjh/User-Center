package register

import (
	"github.com/zjutjh/User-Center/dao/query"
	"github.com/zjutjh/mygo/feishu"
	"github.com/zjutjh/mygo/foundation/kernel"
	"github.com/zjutjh/mygo/ndb"
	"github.com/zjutjh/mygo/nlog"
	"github.com/zjutjh/mygo/usercenter"
)

func Boot() kernel.BootList {
	return kernel.BootList{
		feishu.Boot(),
		// 基础引导器
		nlog.Boot(),

		// Client引导器
		ndb.Boot(),
		//nedis.Boot(),
		// 业务引导器
		usercenter.Boot(),
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
