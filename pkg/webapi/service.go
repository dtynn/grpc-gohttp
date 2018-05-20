package webapi

import (
	"context"
	"errors"
	"net"
	"net/http"
	"time"

	"github.com/dtynn/grpc-gohttp/pkg/logger"
)

// New return new service
func New() *Service {
	return &Service{
		l: logger.Nop(),
	}
}

// Service webapi service
type Service struct {
	GJSONParamer

	mux         http.ServeMux
	middlewares []http.Handler

	server *http.Server
	l      logger.Logger
}

// SetLogger set logger for web api service
func (s *Service) SetLogger(l logger.Logger) {
	if l != nil {
		s.l = l
	}
}

// Use add middlewares
func (s *Service) Use(middlewares ...http.Handler) {
	s.middlewares = append(s.middlewares, middlewares...)
}

// Register register handler
func (s *Service) Register(pattern string, handler http.Handler) {
	s.mux.Handle(pattern, handler)
}

// ServeHTTP implement http.Handler
func (s *Service) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	rw = newResponseWritter(rw)

	for _, m := range s.middlewares {
		if Written(rw) {
			return
		}

		m.ServeHTTP(rw, req)
	}

	if Written(rw) {
		return
	}

	s.mux.ServeHTTP(rw, req)
}

// Run starts the http server
func (s *Service) Run(ctx context.Context, l net.Listener) error {
	if s.server != nil {
		return errors.New("server is running")
	}

	s.server = &http.Server{
		Handler: s,
	}

	errCh := make(chan error, 1)

	go func() {
		err := s.server.Serve(l)
		if err != nil {
			errCh <- err
		}

		close(errCh)
	}()

	select {
	case err := <-errCh:
		return err

	case <-ctx.Done():
		sctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		s.server.Shutdown(sctx)
		return nil
	}
}
