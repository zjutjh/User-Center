package options

import (
	"fmt"
	"net"
	"net/http"
	"time"

	grpcmiddleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpcrecovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/soheilhy/cmux"
	"google.golang.org/grpc"

	"github.com/zjutjh/User-Center-grpc/pkg/apiserver"
	"github.com/zjutjh/User-Center-grpc/pkg/database"
	"github.com/zjutjh/User-Center-grpc/pkg/nacos"
	"github.com/zjutjh/User-Center-grpc/pkg/redis"
	"github.com/zjutjh/User-Center-grpc/pkg/viper"
)

type Options struct {
	ServerRunOptions   *apiserver.ServerRunOptions
	DatabaseRunOptions *database.RunOptions
	RedisRunOptions    *redis.RunOptions
	NacosRunOptions    *nacos.RunOptions
	Debug              bool
}

func NewAPIServerRunOptions() *Options {
	viper.InitViper()
	return &Options{
		ServerRunOptions:   apiserver.NewServerRunOptions(),
		DatabaseRunOptions: database.NewRunOptions(),
		RedisRunOptions:    redis.NewRunOptions(),
		NacosRunOptions:    nacos.NewRunOptions(),
	}
}

func (o *Options) NewAPIServer() (*apiserver.APIServer, error) {
	o.NacosRunOptions.RegisterNacosService()
	o.DatabaseRunOptions.Init()
	o.RedisRunOptions.Init()

	apiServer := &apiserver.APIServer{
		Debug: o.Debug,
	}

	// Create the main listener.
	address := fmt.Sprintf("%s:%d", o.ServerRunOptions.BindAddress, o.ServerRunOptions.InsecurePort)
	l, err := net.Listen("tcp", address)
	if err != nil {
		return nil, err
	}
	// Create a cmux.
	apiServer.CMux = cmux.New(l)

	// Create your protocol servers.
	apiServer.Server = &http.Server{
		Addr:              address,
		ReadHeaderTimeout: 60 * time.Second,
	}

	apiServer.GrpcServer = grpc.NewServer(
		grpc.StreamInterceptor(grpcmiddleware.ChainStreamServer(
			grpcrecovery.StreamServerInterceptor(),
		)),
		grpc.UnaryInterceptor(grpcmiddleware.ChainUnaryServer(
			grpcrecovery.UnaryServerInterceptor(),
		)))

	marshaler := &runtime.JSONPb{}
	marshaler.UseProtoNames = false
	marshaler.EmitUnpopulated = true
	apiServer.GatewayServerMux = runtime.NewServeMux(
		runtime.WithMarshalerOption(runtime.MIMEWildcard, marshaler),
	)

	return apiServer, nil
}
