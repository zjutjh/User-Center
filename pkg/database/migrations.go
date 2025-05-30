package database

import (
	"gorm.io/gorm"

	"github.com/zjutjh/User-Center-grpc/pkg/model"
)

func autoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&model.User{},
		&model.Student{},
	)
}
