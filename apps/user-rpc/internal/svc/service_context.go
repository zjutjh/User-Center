package svc

import (
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zjutjh/User-Center/apps/user-rpc/internal/config"
	"github.com/zjutjh/User-Center/apps/user-rpc/internal/dao/query"
	rpcmodel "github.com/zjutjh/User-Center/apps/user-rpc/internal/model"
	"github.com/zjutjh/User-Center/common/sessionx"
	"gorm.io/gorm"
)

type ServiceContext struct {
	Config  config.Config
	DB      *gorm.DB
	Query   *query.Query
	Session *sessionx.Manager
}

func NewServiceContext(c config.Config) *ServiceContext {
	db, err := rpcmodel.NewDB(c.Mysql)
	logx.Must(err)

	return &ServiceContext{
		Config:  c,
		DB:      db,
		Query:   query.Use(db),
		Session: sessionx.NewManager(c.Session),
	}
}
