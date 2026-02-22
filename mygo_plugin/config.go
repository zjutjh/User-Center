package usercenter

var DefaultConfig = Config{
	Addr: "localhost:8001",
}

type Config struct {
	Addr string `mapstructure:"addr"`
}
