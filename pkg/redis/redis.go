package redis

import (
	"github.com/redis/go-redis/v9"
	"github.com/zjutjh/User-Center-grpc/pkg/viper"
)

var RedisClient *redis.Client

type RunOptions struct {
	Host     string
	Port     string
	DB       int
	Password string
}

func NewRunOptions() *RunOptions {
	Info := &RunOptions{
		Host:     "localhost",
		Port:     "6379",
		DB:       0,
		Password: "",
	}
	if viper.Config.IsSet("redis.host") {
		Info.Host = viper.Config.GetString("redis.host")
	}
	if viper.Config.IsSet("redis.port") {
		Info.Port = viper.Config.GetString("redis.port")
	}
	if viper.Config.IsSet("redis.db") {
		Info.DB = viper.Config.GetInt("redis.db")
	}
	if viper.Config.IsSet("redis.password") {
		Info.Password = viper.Config.GetString("redis.password")
	}
	return Info
}

func (options *RunOptions) Init() {
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     options.Host + ":" + options.Port,
		Password: options.Password,
		DB:       options.DB,
	})
}
