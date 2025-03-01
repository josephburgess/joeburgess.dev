package logging

import (
	"net/http"
	"time"
)

func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		wrw := newResponseWriter(w)
		next.ServeHTTP(wrw, r)
		duration := time.Since(start)

		Log.Infow(
			"HTTP Request",
			"remote_addr", r.RemoteAddr,
			"method", r.Method,
			"uri", r.RequestURI,
			"status", wrw.statusCode,
			"duration", duration,
		)
	})
}

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func newResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{w, http.StatusOK}
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}
