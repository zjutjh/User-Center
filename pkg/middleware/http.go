package middleware

import (
	"log/slog"
	"net/http"
	"strings"
	"time"
)

type ResponseWriter struct {
	http.ResponseWriter
	StatusCode int
}

func NewResponseWriter(w http.ResponseWriter) *ResponseWriter {
	return &ResponseWriter{
		ResponseWriter: w,
		StatusCode:     http.StatusOK,
	}
}

func (r *ResponseWriter) WriteHeader(code int) {
	r.StatusCode = code
	r.ResponseWriter.WriteHeader(code)
}

func (r *ResponseWriter) Flush() {
	r.ResponseWriter.(http.Flusher).Flush()
}

func LogRequestAndResponse(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		ws := NewResponseWriter(w)
		// The Grpc-Metadata-X-Real-IP header is added to get the Real IP and port in the grpc middleware
		IPAddress := ReadHTTPIP(r)
		RequestURI := ReadXRequestURI(r)
		r.Header.Set("Grpc-Metadata-X-Real-IP", IPAddress)
		r.Header.Set("Grpc-Metadata-X-Request-URI", RequestURI)
		r.Header.Set("Grpc-Metadata-X-Request-URL", r.URL.Path)
		handler.ServeHTTP(ws, r)
		slog.Info("[kantaloupe]", "time",
			time.Now().Format("2006/01/02 - 15:04:05"),
			"code", ws.StatusCode, "cost",
			time.Since(start).String(), "ip",
			IPAddress,
			"method", r.Method, "uri",
			r.URL.Path)
	})
}

func ReadHTTPIP(r *http.Request) string {
	IPAddress := r.Header.Get("X-Forwarded-For")
	if IPAddress == "" {
		IPAddress = r.Header.Get("X-Real-Ip")
	}
	if IPAddress == "" {
		IPAddress = r.RemoteAddr
	}
	return IPAddress
}

func ReadXRequestURI(r *http.Request) string {
	header := r.Header.Get("X-Forwarded-Client-Cert")
	if header == "" {
		return ""
	}
	for _, v := range strings.Split(header, ";") {
		if strings.HasPrefix(v, "URI=") {
			return v
		}
	}
	return ""
}

const (
	green   = "\033[97;42m"
	white   = "\033[90;47m"
	yellow  = "\033[90;43m"
	red     = "\033[97;41m"
	blue    = "\033[97;44m"
	magenta = "\033[97;45m"
	cyan    = "\033[97;46m"
	reset   = "\033[0m"

	GRPCHeaderClusterKey = "x-cluster"

	HTTPHeaderClusterKey = "Grpc-Metadata-X-Cluster"
)

// HTTPCodeColor is the ANSI color for appropriately logging http status code to a terminal.
func HTTPCodeColor(code int) string {
	switch {
	case code >= http.StatusOK && code < http.StatusMultipleChoices:
		return green
	case code >= http.StatusMultipleChoices && code < http.StatusBadRequest:
		return white
	case code >= http.StatusBadRequest && code < http.StatusInternalServerError:
		return yellow
	default:
		return red
	}
}

// MethodColor is the ANSI color for appropriately logging http method to a terminal.
func MethodColor(method string) string {
	switch method {
	case http.MethodGet:
		return blue
	case http.MethodPost:
		return cyan
	case http.MethodPut:
		return yellow
	case http.MethodDelete:
		return red
	case http.MethodPatch:
		return green
	case http.MethodHead:
		return magenta
	case http.MethodOptions:
		return white
	default:
		return reset
	}
}

// ResetColor resets all escape attributes.
func ResetColor() string {
	return reset
}

func ConvertHTTPParams(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ws := NewResponseWriter(w)

		// The Grpc-Metadata-X-Real-IP header is added to get the Real IP and port in the grpc middleware
		IPAddress := ReadHTTPIP(r)
		RequestURI := ReadXRequestURI(r)

		r.Header.Set("Grpc-Metadata-X-Real-IP", IPAddress)
		r.Header.Set("Grpc-Metadata-X-Request-URI", RequestURI)
		handler.ServeHTTP(ws, r)
	})
}
