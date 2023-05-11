package http

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/labstack/echo/v4"
	"github.com/webws/go-moda/logger"
	"github.com/webws/go-moda/transport"
)

var _ transport.Server = (*Server)(nil)

const PprofPrefix = "/debug/pprof"

// default http address
const address = ":8081"

type Server struct {
	ctx         context.Context
	network     string
	address     string
	handle      HTTPServer
	pprofPrefix string
	tracing     bool
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

// tracing with server tracing.
func WithTracing(tracing bool) ServerOption {
	return func(s *Server) {
		s.tracing = tracing
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
		logger.Infow("[HTTP] server address is empty, use default address", "address", address)
		srv.address = ":8081"
	}
	srv.handle.PprofRegister(srv.pprofPrefix)
	if srv.tracing {
		// 启用链路追踪
		srv.handle.EnableTracing()
	}
	return srv
}

func NewEchoHttpServer(opts ...ServerOption) (*echo.Echo, *Server) {
	echoServer := newEchoServer()
	opts = append(opts, WitchHandle(echoServer))
	srv := NewServer(opts...)
	return echoServer.GetServer(), srv
}

func NewGinHttpServer(opts ...ServerOption) (*gin.Engine, *Server) {
	ginServer := newGinServer()
	opts = append(opts, WitchHandle(ginServer))
	srv := NewServer(opts...)
	return ginServer.GetServer(), srv
}

func NewNetHttpServer(opts ...ServerOption) (*http.ServeMux, *Server) {
	netHttpServer := newNetHTTPServer()
	opts = append(opts, WitchHandle(netHttpServer))
	srv := NewServer(opts...)
	return netHttpServer.GetServer(), srv
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
	if err := s.handle.Stop(ctx); err != nil {
		logger.Errorw("[HTTP] server stop failed", "listen_addr", s.address, "err", err)
		return err
	}
	logger.Infow("[HTTP] server stopped")
	return nil
}
