package middleware

import (
	"net/http"
	"time"

	"go.uber.org/zap"
)

type RequestLogger struct {
	logger *zap.Logger
}

func NewRequestLogger(logger *zap.Logger) RequestLogger {
	return RequestLogger{logger: logger}
}

type wrappedWriter struct {
	http.ResponseWriter
	status int
	len    int
}

func (w *wrappedWriter) WriteHeader(status int) {
	w.ResponseWriter.WriteHeader(status)
	w.status = status
}

func (w *wrappedWriter) Write(b []byte) (len int, err error) {
	len, err = w.ResponseWriter.Write(b)
	w.len += len
	return
}

func (l RequestLogger) Wrap(callback http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		writer := &wrappedWriter{w, 0, 0}
		start := time.Now()
		callback(writer, r)
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
	}
}
