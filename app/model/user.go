package model

import (
	"encoding/json"
	"time"
)

type User struct {
	StudentId    string
	UserId       int `gorm:"primary_key;AUTO_INCREMENT"`
	Password     string
	Email        string
	Type         uint8           // 0: 本科生 1: 研究生
	BoundSystems json.RawMessage `gorm:"type:json"` //  绑定的系统
	CreateTime   time.Time
}
