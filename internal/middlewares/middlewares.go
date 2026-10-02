package middlewares

import (
	"log"
	"net/http"
	"time"
)

type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (lw *loggingResponseWriter) WriteHeader(statusCode int) {
	lw.statusCode = statusCode
	lw.ResponseWriter.WriteHeader(statusCode)
}

func (lw *loggingResponseWriter) Write(data []byte) (int, error) {
	if lw.statusCode == 0 {
		lw.statusCode = http.StatusOK
	}
	return lw.ResponseWriter.Write(data)
}

func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		
		lw := &loggingResponseWriter{
			ResponseWriter: w,
			statusCode: 0,
		}

		next.ServeHTTP(lw, r)
		log.Printf("[LOG] %s %s status=%d duration=%v", r.Method, r.URL.Path, lw.statusCode, time.Since(start))
	})
}