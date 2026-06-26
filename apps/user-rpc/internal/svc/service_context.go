package svc

import (
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zjutjh/User-Center/apps/user-rpc/internal/config"
	"github.com/zjutjh/User-Center/apps/user-rpc/internal/dao/query"
	"github.com/zjutjh/User-Center/apps/user-rpc/internal/dao/repo"
	"github.com/zjutjh/User-Center/apps/user-rpc/internal/domain/credential"
	"github.com/zjutjh/User-Center/apps/user-rpc/internal/domain/oauth"
	"github.com/zjutjh/User-Center/apps/user-rpc/internal/infra/mysql"
	"github.com/zjutjh/User-Center/common/sessionx"
	"gorm.io/gorm"
)

type ServiceContext struct {
	Config      config.Config
	DB          *gorm.DB
	Query       *query.Query
	UserRepo    *repo.UserRepo
	StudentRepo *repo.StudentRepo
	Session     *sessionx.Manager
	Credential  *credential.Service
	OAuth       *oauth.Service
}

func NewServiceContext(c config.Config) *ServiceContext {
	db, err := mysql.NewDB(c.Mysql)
	logx.Must(err)

	credentialSvc, err := credential.NewService(c.Credential.SecretKey)
	logx.Must(err)

	q := query.Use(db)

	return &ServiceContext{
		Config:      c,
		DB:          db,
		Query:       q,
		UserRepo:    repo.NewUserRepo(q),
		StudentRepo: repo.NewStudentRepo(q),
		Session:     sessionx.NewManager(c.Session),
		Credential:  credentialSvc,
		OAuth:       oauth.NewService(),
	}
}
