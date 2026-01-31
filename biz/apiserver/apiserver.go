package apiserver

import (
	"context"
	"log/slog"
	"net"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/zjutjh/User-Center/biz/handler"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	serverapi "github.com/zjutjh/User-Center/api/user/v1alpha1"
	"github.com/zjutjh/User-Center/biz/middleware"
)

type APIServer struct {
	Debug            bool
	Server           *http.Server
	GrpcServer       *grpc.Server
	GatewayServerMux *runtime.ServeMux
	router           *mux.Router

	GrpcListener net.Listener
	HttpListener net.Listener
}

func (s *APIServer) PrepareRun(ctx context.Context) error {
	s.router = mux.NewRouter()
	if err := s.registerGrpcServices(ctx); err != nil {
		return err
	}
	// mux middleware
	s.router.Use(middleware.LogRequestAndResponse)
	s.registerHTTPAPIs()
	s.Server.Handler = s.router
	return nil
}

func (s *APIServer) registerGrpcServices(ctx context.Context) error {
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}
	serverapi.RegisterUserCenterServiceServer(s.GrpcServer, handler.NewUserHandler())
	if err := serverapi.RegisterUserCenterServiceHandlerFromEndpoint(ctx, s.GatewayServerMux, s.GrpcListener.Addr().String(), opts); err != nil {
		return err
	}
	s.router.PathPrefix("/api").Handler(s.GatewayServerMux)
	return nil
}

func (s *APIServer) registerHTTPAPIs() {
	healthRouter := s.router.PathPrefix("/").Subrouter()
	healthRouter.HandleFunc("/healthz", livenessProbe)
	healthRouter.HandleFunc("/readyz", readinessProbe)
}

func (s *APIServer) Run(ctx context.Context) error {
	go func() {
		err := s.GrpcServer.Serve(s.GrpcListener)
		if err != nil {
			slog.Error("Failed to start grpc server", err)
		}
	}()
	go func() {
		err := s.Server.Serve(s.HttpListener)
		if err != nil {
			slog.Error("Failed to start http server", err)
		}
	}()
	slog.Info("Serving...")
	<-ctx.Done()
	return nil
}

func livenessProbe(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("ok"))
}

func readinessProbe(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("ok"))
}
