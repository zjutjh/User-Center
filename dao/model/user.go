package model

import (
	"time"
)

type User struct {
	StudentId  string
	UserId     int `gorm:"primary_key;AUTO_INCREMENT"`
	Password   string
	Email      string
	CreateTime time.Time
}
