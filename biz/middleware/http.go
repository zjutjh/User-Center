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
