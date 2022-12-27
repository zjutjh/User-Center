package configStart

import (
	"usercenter/config/config"
	"usercenter/config/database"
	"usercenter/config/email"
	"usercenter/config/redis"
)

func Init() {
	config.InitConfig()
	database.Init()
	email.Init()
	redis.Init()
}
