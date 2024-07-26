package handler

import (
	"compress/gzip"
	"net/http"
	"strings"
)

type Gzip struct{}

func (g *Gzip) GzipMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			ww := NewWrappedGzipResponseWriter(rw)
			ww.Header().Set("Content-Encoding", "gzip")
			defer ww.Flush()
			next.ServeHTTP(ww, r)
			return
		}
		next.ServeHTTP(rw, r)
	})
}

type WrappedGzipResponseWriter struct {
	rw http.ResponseWriter
	gw *gzip.Writer
}

func NewWrappedGzipResponseWriter(rw http.ResponseWriter) *WrappedGzipResponseWriter {
	gw := gzip.NewWriter(rw)
	return &WrappedGzipResponseWriter{rw, gw}
}

func (ww *WrappedGzipResponseWriter) Header() http.Header {
	return ww.rw.Header()
}

func (ww *WrappedGzipResponseWriter) Write(d []byte) (int, error) {
	return ww.gw.Write(d)
}

func (ww *WrappedGzipResponseWriter) WriteHeader(statusCode int) {
	ww.rw.WriteHeader(statusCode)
}
func (ww *WrappedGzipResponseWriter) Flush() {
	ww.gw.Flush()
	ww.gw.Close()
}
