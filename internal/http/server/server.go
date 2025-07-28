// server/server.go (as provided, with minor clarification)
package server

import (
	"log/slog"
	"net/http"
	"sync"
)

type HTTPHandler interface {
	ServeHTTP(http.ResponseWriter, *http.Request)
}

type routes map[string]HTTPHandler

type server struct {
	routes routes
	logger *slog.Logger
	init   sync.Once
	mux    *http.ServeMux
}

func New(r routes, l *slog.Logger) *server {
	logger := l.WithGroup("server")
	return &server{
		routes: r,
		logger: logger,
		mux:    http.NewServeMux(),
	}
}

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func loggingMiddleware(next HTTPHandler, logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rw := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}
		next.ServeHTTP(rw, r)
		logger.Info("HTTP request",
			"method", r.Method,
			"path", r.URL.Path,
			"status", rw.statusCode,
			"remote_addr", r.RemoteAddr,
			"user_agent", r.UserAgent(),
		)
	}
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.init.Do(func() {
		for route, handler := range s.routes {
			s.mux.HandleFunc(route, loggingMiddleware(handler, s.logger))
		}
	})
	s.mux.ServeHTTP(w, r)
}
