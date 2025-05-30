package apiserver

import (
	"context"
	"log/slog"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/soheilhy/cmux"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	serverapi "github.com/zjutjh/User-Center-grpc/api/v1"
	"github.com/zjutjh/User-Center-grpc/pkg/apiserver/bff"
	"github.com/zjutjh/User-Center-grpc/pkg/middleware"
)

type APIServer struct {
	stopCh           chan struct{}
	Debug            bool
	Server           *http.Server
	GrpcServer       *grpc.Server
	GatewayServerMux *runtime.ServeMux
	CMux             cmux.CMux
	router           *mux.Router
	PrometheusAddr   string
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
	s.stopCh = make(chan struct{})

	return nil
}

func (s *APIServer) registerGrpcServices(ctx context.Context) error {
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	serverapi.RegisterUserServer(s.GrpcServer, bff.NewUserHandler())
	err := serverapi.RegisterUserHandlerFromEndpoint(ctx, s.GatewayServerMux, s.Server.Addr, opts)
	if err != nil {
		return err
	}

	s.router.PathPrefix("/api").Handler(s.GatewayServerMux)
	return err
}

func (s *APIServer) registerHTTPAPIs() {
	healthRouter := s.router.PathPrefix("/").Subrouter()
	healthRouter.HandleFunc("/healthz", livenessProbe)
	healthRouter.HandleFunc("/readyz", readinessProbe)
}

func (s *APIServer) printRouters() error {
	return s.router.Walk(func(route *mux.Route, _ *mux.Router, _ []*mux.Route) error {
		pathTemplate, err := route.GetPathTemplate()
		if err == nil {
			slog.Info("ROUTE:", "template", pathTemplate)
		}
		pathRegexp, err := route.GetPathRegexp()
		if err == nil {
			slog.Info("Path regexp:", "regexp", pathRegexp)
		}
		queriesTemplates, err := route.GetQueriesTemplates()
		if err == nil {
			slog.Info("Queries templates:", "templates", strings.Join(queriesTemplates, ","))
		}
		queriesRegexps, err := route.GetQueriesRegexp()
		if err == nil {
			slog.Info("Queries regexps:", "regexps", strings.Join(queriesRegexps, ","))
		}
		methods, err := route.GetMethods()
		if err == nil {
			slog.Info("Methods:", "methods", strings.Join(methods, ","))
		}
		return nil
	})
}

func (s *APIServer) Run(ctx context.Context) error {
	s.waitForResourceSync(ctx)

	// Match connections in order:
	// First grpc, then HTTP, and otherwise Go RPC/TCP.
	grpcListener := s.CMux.Match(cmux.HTTP2())
	httpListener := s.CMux.Match(cmux.HTTP1Fast("PATCH"))

	// Use the muxed listeners for your servers.
	go func() {
		err := s.GrpcServer.Serve(grpcListener)
		if err != nil {
			slog.Error("Failed to start grpc server", err)
		}
	}()

	go func() {
		err := s.Server.Serve(httpListener)
		if err != nil {
			slog.Error("Failed to start http server", err)
		}
	}()

	// Start serving!
	slog.Info("Serving...")
	return s.CMux.Serve()
}

func (s *APIServer) waitForResourceSync(_ context.Context) {
	slog.Info("Start cache objects")

	slog.Info("Finished caching objects")
}

func livenessProbe(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("ok"))
}

func readinessProbe(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("ok"))
}
