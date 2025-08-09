package middleware

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"
)

var compressingEnableEncodings = []string{
	"application/json",
	"text/html",
}

type compressedWriter struct {
	httpWriter http.ResponseWriter
	writer     io.Writer
	status     int
}

func newCompressedWriter(httpWriter http.ResponseWriter) *compressedWriter {
	return &compressedWriter{
		httpWriter: httpWriter,
		status:     http.StatusOK,
	}
}

func (c *compressedWriter) Header() http.Header {
	return c.httpWriter.Header()
}

func (c *compressedWriter) WriteHeader(statusCode int) {
	c.status = statusCode
	// headers will be written later
}

func (c *compressedWriter) Write(p []byte) (int, error) {
	c.initWriter()
	return c.writer.Write(p)
}

func (c *compressedWriter) initWriter() {
	if c.writer != nil {
		return
	}

	defer c.httpWriter.WriteHeader(c.status)

	contentEncoding := c.Header().Get("Content-Type")
	enableCompression := false
	for _, compressingEnableEncoding := range compressingEnableEncodings {
		enableCompression = enableCompression || contentEncoding == compressingEnableEncoding
	}

	if !enableCompression {
		c.writer = c.httpWriter
		return
	}

	c.Header().Set("Content-Encoding", "gzip")
	c.writer = gzip.NewWriter(c.httpWriter)
}

func (c *compressedWriter) Close() error {
	if c.writer == nil {
		c.httpWriter.WriteHeader(c.status)
		return nil
	}

	if closer, ok := c.writer.(io.Closer); ok {
		return closer.Close()
	}

	return nil
}

func Compression(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		acceptEncodings := strings.Split(r.Header.Get("Accept-Encoding"), ",")
		supportsGzip := false
		for _, acceptEncoding := range acceptEncodings {
			supportsGzip = supportsGzip || strings.HasPrefix(strings.TrimSpace(acceptEncoding), "gzip")
		}

		if !supportsGzip {
			handler(w, r)
			return
		}

		cw := newCompressedWriter(w)
		defer cw.Close()
		handler(cw, r)
	}
}
