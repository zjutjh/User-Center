package model

import "time"

// 系统绑定信息的结构体，用于存储 JSON 数据
type SystemBinding struct {
	SystemName SystemNameEnum `json:"SystemName"`
	BoundAt    time.Time      `json:"BoundAt"`
}

type SystemNameEnum uint

const (
	wjh  SystemNameEnum = 0
	foru SystemNameEnum = 1
)
