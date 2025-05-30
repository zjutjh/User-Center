package apiserver

import "github.com/zjutjh/User-Center-grpc/pkg/viper"

type ServerRunOptions struct {
	// server bind address
	BindAddress string

	// insecure port number
	InsecurePort int

	// secure port number
	SecurePort int

	// tls cert file
	TLSCertFile string

	// tls private key file
	TLSPrivateKey string

	// insecure grpc port number
	InsecureGRPCPort int
}

func NewServerRunOptions() *ServerRunOptions {
	// create default server run options
	Info := ServerRunOptions{
		BindAddress:      "0.0.0.0",
		InsecurePort:     8000,
		InsecureGRPCPort: 8001,
		SecurePort:       0,
		TLSCertFile:      "",
		TLSPrivateKey:    "",
	}

	if viper.Config.IsSet("server.address") {
		Info.BindAddress = viper.Config.GetString("server.address")
	}
	if viper.Config.IsSet("server.port.insecure") {
		Info.InsecurePort = viper.Config.GetInt("server.port.insecure")
	}
	if viper.Config.IsSet("server.port.insecureGRPC") {
		Info.InsecureGRPCPort = viper.Config.GetInt("server.port.insecureGRPC")
	}
	if viper.Config.IsSet("server.port.secure") {
		Info.SecurePort = viper.Config.GetInt("server.port.secure")
	}
	if viper.Config.IsSet("server.tls.certFile") {
		Info.TLSCertFile = viper.Config.GetString("server.tls.certFile")
	}
	if viper.Config.IsSet("server.tls.privateKey") {
		Info.TLSPrivateKey = viper.Config.GetString("server.tls.privateKey")
	}
	return &Info
}
