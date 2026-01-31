package apiserver

import "github.com/zjutjh/User-Center/biz/viper"

type ServerRunOptions struct {
	// server bind address
	BindAddress string
	// insecure grpc port number
	GRPCPort int
	HttpPort int
}

func NewServerRunOptions() *ServerRunOptions {
	// create default server run options
	Info := ServerRunOptions{
		BindAddress: "0.0.0.0",
		GRPCPort:    8000,
		HttpPort:    8001,
	}

	if viper.Config.IsSet("server.address") {
		Info.BindAddress = viper.Config.GetString("server.address")
	}
	if viper.Config.IsSet("server.port.grpc") {
		Info.GRPCPort = viper.Config.GetInt("server.port.grpc")
	}
	if viper.Config.IsSet("server.port.http") {
		Info.HttpPort = viper.Config.GetInt("server.port.http")
	}
	return &Info
}
