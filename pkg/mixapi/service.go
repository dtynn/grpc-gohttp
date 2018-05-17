package mixapi

import (
	"context"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/dtynn/grpc-gohttp/pkg/grpcapi"
	"github.com/dtynn/grpc-gohttp/pkg/logger"
	"github.com/dtynn/grpc-gohttp/pkg/webapi"

	"github.com/cockroachdb/cmux"
)

// Mix return a mixed service
func Mix(g *grpcapi.Service, w *webapi.Service) *Service {
	return &Service{
		g: g,
		w: w,
		l: logger.Nop(),
	}
}

// Service mixed service
type Service struct {
	g *grpcapi.Service
	w *webapi.Service
	l logger.Logger
}

// SetLogger set logger for service
func (s *Service) SetLogger(l logger.Logger) {
	if l != nil {
		s.l = l
	}

	if s.g != nil {
		s.g.SetLogger(l)
	}

	if s.w != nil {
		s.w.SetLogger(l)
	}
}

// Web return webapi service
func (s *Service) Web() *webapi.Service {
	return s.w
}

// GRPC return grpc service
func (s *Service) GRPC() *grpcapi.Service {
	return s.g
}

// Listen start a listener
func (s *Service) Listen(addr string) {
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		s.l.Fatal(err)
	}

	s.l.Info("listen on", addr)

	m := cmux.New(lis)
	gL := m.Match(cmux.HTTP2HeaderField("content-type", "application/grpc"))
	wL := m.Match(cmux.Any())

	errCh := make(chan error, 1)
	go func() {
		if err := m.Serve(); err != nil {
			errCh <- err
		}

		close(errCh)
	}()

	ctx, cancel := context.WithCancel(context.Background())
	var wg sync.WaitGroup

	if s.g != nil {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := s.g.Run(ctx, gL); err != nil {
				s.l.Fatal("grpc api service error", err)
			}
		}()
	}

	if s.w != nil {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := s.w.Run(ctx, wL); err != nil {
				s.l.Fatal("web api service error", err)
			}
		}()
	}

	sigCh := make(chan os.Signal, 3)

	go func() {
		signal.Notify(sigCh, syscall.SIGTERM, syscall.SIGINT, syscall.SIGABRT)
	}()

	select {
	case err := <-errCh:
		s.l.Warn("cmux serve", err)

	case sig := <-sigCh:
		s.l.Info("captured signal", sig)
	}

	cancel()

	wg.Wait()
	s.l.Info("server stopped")
}
