package model

import "time"

type User struct {
	SID        int
	Password   string
	Email      string
	CreateTime time.Time
}
