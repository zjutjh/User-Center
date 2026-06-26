package svc

import (
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zjutjh/User-Center/apps/user-rpc/internal/config"
	"github.com/zjutjh/User-Center/apps/user-rpc/internal/dao/query"
	"github.com/zjutjh/User-Center/apps/user-rpc/internal/domain/credential"
	"github.com/zjutjh/User-Center/apps/user-rpc/internal/infra/mysql"
	"github.com/zjutjh/User-Center/common/sessionx"
	"gorm.io/gorm"
)

type ServiceContext struct {
	Config     config.Config
	DB         *gorm.DB
	Query      *query.Query
	Session    *sessionx.Manager
	Credential *credential.Service
}

func NewServiceContext(c config.Config) *ServiceContext {
	db, err := mysql.NewDB(c.Mysql)
	logx.Must(err)

	credentialSvc, err := credential.NewService(c.Credential.SecretKey)
	logx.Must(err)

	return &ServiceContext{
		Config:     c,
		DB:         db,
		Query:      query.Use(db),
		Session:    sessionx.NewManager(c.Session),
		Credential: credentialSvc,
	}
}
