package grpcserver

import (
	"context"
	"net"

	"google.golang.org/grpc"
)

type Server struct {
	GRPC   *grpc.Server
	ln     net.Listener
	closed chan struct{}
}

type Options struct {
	Addr   string
	Server *grpc.Server
}

func New(opts Options) (*Server, error) {
	ln, err := net.Listen("tcp", opts.Addr)
	if err != nil {
		return nil, err
	}
	s := opts.Server
	if s == nil {
		s = grpc.NewServer()
	}
	return &Server{GRPC: s, ln: ln, closed: make(chan struct{})}, nil
}

func (s *Server) Start() {
	go func() {
		_ = s.GRPC.Serve(s.ln)
		close(s.closed)
	}()
}

func (s *Server) Shutdown(ctx context.Context) error {
	done := make(chan struct{})
	go func() { s.GRPC.GracefulStop(); close(done) }()
	select {
	case <-done:
		<-s.closed
		return nil
	case <-ctx.Done():
		s.GRPC.Stop()
		<-s.closed
		return ctx.Err()
	}
}
