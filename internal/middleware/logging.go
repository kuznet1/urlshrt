package middleware

import (
	"net/http"
	"time"

	"go.uber.org/zap"
)

var Log *zap.Logger

func init() {
	var err error
	Log, err = zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
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

func Logging(callback http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		writer := &wrappedWriter{w, 0, 0}
		start := time.Now()
		callback(writer, r)
		elapsed := time.Since(start)
		Log.Info("request:",
			zap.String("uri", r.RequestURI),
			zap.String("method", r.Method),
			zap.Duration("time", elapsed),
		)
		Log.Info("response:",
			zap.Int("status", writer.status),
			zap.Int("size", writer.len),
		)
	}
}
