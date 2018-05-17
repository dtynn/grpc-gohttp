package grpcapi

import (
	"context"
	"net"

	"github.com/dtynn/grpc-gohttp/pkg/logger"
	"google.golang.org/grpc"
)

// New returns a grpc service
func New(opt ...grpc.ServerOption) *Service {
	return &Service{
		Server: grpc.NewServer(opt...),
		l:      logger.Nop(),
	}
}

// Service grpc api service
type Service struct {
	*grpc.Server
	l logger.Logger
}

// SetLogger set logger for grpc api service
func (s *Service) SetLogger(l logger.Logger) {
	if l != nil {
		s.l = l
	}
}

// Run starts the grpc server
func (s *Service) Run(ctx context.Context, l net.Listener) error {
	errCh := make(chan error, 1)

	go func() {
		err := s.Server.Serve(l)
		if err != nil {
			errCh <- err
		}

		close(errCh)
	}()

	select {
	case err := <-errCh:
		return err

	case <-ctx.Done():
		s.Server.GracefulStop()
		return nil
	}
}
