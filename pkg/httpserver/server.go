package httpserver

import (
	"context"
	"net"
	"net/http"
	"time"
)

type Server struct {
	http   *http.Server
	ln     net.Listener
	closed chan struct{}
}

type Options struct {
	Addr         string
	Handler      http.Handler
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
}

func New(opts Options) (*Server, error) {
	s := &http.Server{
		Addr:         opts.Addr,
		Handler:      opts.Handler,
		ReadTimeout:  opts.ReadTimeout,
		WriteTimeout: opts.WriteTimeout,
		IdleTimeout:  opts.IdleTimeout,
	}

	ln, err := net.Listen("tcp", opts.Addr)
	if err != nil {
		return nil, err
	}

	return &Server{
		http:   s,
		ln:     ln,
		closed: make(chan struct{}),
	}, nil
}

func (s *Server) Start() {
	go func() {
		_ = s.http.Serve(s.ln)
		close(s.closed)
	}()
}

func (s *Server) Shutdown(ctx context.Context) error {
	defer func() { <-s.closed }()
	return s.http.Shutdown(ctx)
}
