package nacos

import (
	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
	"github.com/zjutjh/User-Center-grpc/pkg/viper"
	"log"
)

type RunOptions struct {
	User string
	Pass string
	Port uint64
	Host string
	Name string
}

func NewRunOptions() *RunOptions {
	Info := &RunOptions{
		User: "nacos",
		Pass: "nacos",
		Port: 8848,
		Host: "127.0.0.1",
		Name: "user-center-grpc",
	}
	if viper.Config.IsSet("server.name") {
		Info.Name = viper.Config.GetString("server.name")
	}
	if viper.Config.IsSet("nacos.host") {
		Info.Host = viper.Config.GetString("nacos.host")
	}
	if viper.Config.IsSet("nacos.port") {
		Info.Port = viper.Config.GetUint64("nacos.port")
	}
	if viper.Config.IsSet("nacos.user") {
		Info.User = viper.Config.GetString("nacos.user")
	}
	if viper.Config.IsSet("nacos.pass") {
		Info.Pass = viper.Config.GetString("nacos.pass")
	}
	return Info
}

func (options *RunOptions) RegisterNacosService() {
	serverConfigs := []constant.ServerConfig{
		*constant.NewServerConfig(options.Host, options.Port),
	}
	clientConfig := *constant.NewClientConfig(
		constant.WithNamespaceId("public"),
		constant.WithTimeoutMs(5000),
		constant.WithNotLoadCacheAtStart(true),
		constant.WithLogDir("./nacos/log"),
		constant.WithCacheDir("./nacos/cache"),
		constant.WithLogLevel("info"),
		constant.WithAppName(options.Name),
		constant.WithUsername(options.User),
		constant.WithPassword(options.Pass),
	)
	namingClient, err := clients.NewNamingClient(
		vo.NacosClientParam{
			ClientConfig:  &clientConfig,
			ServerConfigs: serverConfigs,
		},
	)
	if err != nil {
		log.Println("Nacos client init failed:", err)
		return
	}
	_, err = namingClient.RegisterInstance(vo.RegisterInstanceParam{
		Ip:          options.Host, // 动态传入
		Port:        options.Port, // 动态传入
		ServiceName: options.Name, // 动态传入
		Weight:      10,
		Enable:      true,
		Healthy:     true,
		Ephemeral:   true,
	})
	if err != nil {
		log.Println("Nacos service registration failed:", err)
		return
	}
	log.Printf("Nacos service %s registered successfully at %s:%d\n", options.Name, options.Host, options.Port)
}
