package middleware

import (
	"net/http"
	"time"

	"go.uber.org/zap"
)

// RequestLogger is a public struct of the package. It exposes the core data for this project.
type RequestLogger struct {
	logger *zap.Logger
}

// NewRequestLogger performs a public package operation. Top-level handler/function.
func NewRequestLogger(logger *zap.Logger) RequestLogger {
	return RequestLogger{logger: logger}
}

type wrappedWriter struct {
	http.ResponseWriter
	status int
	len    int
}

// WriteHeader is a method that provides public behavior for the corresponding type.
func (w *wrappedWriter) WriteHeader(status int) {
	w.ResponseWriter.WriteHeader(status)
	w.status = status
}

// Write is a method that provides public behavior for the corresponding type.
func (w *wrappedWriter) Write(b []byte) (len int, err error) {
	len, err = w.ResponseWriter.Write(b)
	w.len += len
	return
}

// Logging is a method that provides public behavior for the corresponding type.
func (l RequestLogger) Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		writer := &wrappedWriter{w, 0, 0}
		start := time.Now()
		next.ServeHTTP(writer, r)
		elapsed := time.Since(start)
		l.logger.Info("request:",
			zap.String("uri", r.RequestURI),
			zap.String("method", r.Method),
			zap.Duration("time", elapsed),
		)
		l.logger.Info("response:",
			zap.Int("status", writer.status),
			zap.Int("size", writer.len),
		)
	})
}
