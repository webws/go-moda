package http

import (
	"context"
	"net/http"

	"github.com/webws/go-moda/logger"
)

const PprofPrefix = "/debug/pprof"

type Server struct {
	ctx         context.Context
	network     string
	address     string
	handle      HTTPServer
	pprofPrefix string
}

// ServerOption is an HTTP server option.
type ServerOption func(*Server)

// Address with server address.
func WithAddress(addr string) ServerOption {
	return func(s *Server) {
		s.address = addr
	}
}

// Handle with HTTP server handler.
func WitchHandle(h HTTPServer) ServerOption {
	return func(s *Server) {
		s.handle = h
	}
}

// NewServer creates an HTTP server by options.
func NewServer(opts ...ServerOption) *Server {
	srv := &Server{
		network:     "tcp",
		pprofPrefix: PprofPrefix,
	}
	for _, o := range opts {
		o(srv)
	}
	if srv.address == "" {
		srv.address = ":8081"
	}
	srv.handle.PprofRegister(srv.pprofPrefix)
	return srv
}

// Start start the HTTP server.
func (s *Server) Start(ctx context.Context) error {
	s.ctx = ctx
	// logger.Infow("[HTTP] server started", "listen_addr", s.address)
	if err := s.handle.Start(s.address); err != nil && err != http.ErrServerClosed {
		logger.Errorw("[HTTP] server start failed", "listen_addr", s.address, "err", err)
		return err
	}

	return nil
}

// stop stop the HTTP server.
func (s *Server) Stop(ctx context.Context) error {
	ctx = context.Background()
	if err := s.handle.Stop(ctx); err != nil {
		logger.Errorw("[HTTP] server stop failed", "listen_addr", s.address, "err", err)
		return err
	}
	logger.Infow("[HTTP] server stopped")
	return nil
}
