package database

import (
	"gorm.io/gorm"
	"usercenter/app/model"
)

func autoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&model.User{},
		&model.Student{},
	)
}
