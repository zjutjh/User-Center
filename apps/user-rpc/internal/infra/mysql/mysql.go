package mysql

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/zjutjh/User-Center/common/dbx"
)

func NewDB(conf dbx.MysqlConf) (*gorm.DB, error) {
	return gorm.Open(mysql.Open(conf.DSN()), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
	})
}
