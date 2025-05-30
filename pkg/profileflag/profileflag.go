package profileflag

import (
	"log/slog"
	"net"
	"net/http"
	"net/http/pprof"
	"os"
	"strconv"
	"time"

	"github.com/spf13/pflag"
)

const (
	// ReadHeaderTimeout is the amount of time allowed to read
	// request headers.
	// HTTP timeouts are necessary to expire inactive connections
	// and failing to do so might make the application vulnerable
	// to attacks like slowloris which work by sending data very slow,
	// which in case of no timeout will keep the connection active
	// eventually leading to a denial-of-service (DoS) attack.
	// References:
	// - https://en.wikipedia.org/wiki/Slowloris_(computer_security)
	ReadHeaderTimeout = 32 * time.Second
)

// Options are options for pprof.
type Options struct {
	// ProfilingBindAddress is the IP address for the profiling server to serve on.
	ProfilingBindAddress string

	// ProfilingPort is the port of the localhost profiling endpoint (set to 0 to disable)
	ProfilingPort int32
}

// AddFlags adds flags to the specified FlagSet.
func (o *Options) AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(&o.ProfilingBindAddress, "profiling-bind-address", o.ProfilingBindAddress,
		"The IP address for the profiling server to serve on "+
			"(set to '0.0.0.0' or '::' for listening in all interfaces and IP families)")
	fs.Int32Var(&o.ProfilingPort, "profiling-port", o.ProfilingPort,
		"The port of the localhost profiling endpoint (set to 0 to disable)")
}

func installHandlerForPProf(mux *http.ServeMux) {
	mux.HandleFunc("/debug/pprof/", pprof.Index)
	mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	mux.HandleFunc("/debug/pprof/trace", pprof.Trace)
}

// ListenAndServe start a http server to enable pprof.
func ListenAndServe(opts Options) {
	if opts.ProfilingPort > 0 {
		mux := http.NewServeMux()
		installHandlerForPProf(mux)
		addr := net.JoinHostPort(opts.ProfilingBindAddress, strconv.Itoa(int(opts.ProfilingPort)))
		slog.Info("Starting profiling on address", "addr", addr)
		go func() {
			httpServer := http.Server{
				Addr:              addr,
				Handler:           mux,
				ReadHeaderTimeout: ReadHeaderTimeout,
			}
			if err := httpServer.ListenAndServe(); err != nil {
				slog.Error("Failed to start profiling server", err)
				os.Exit(1)
			}
		}()
	}
}
