package grpc

import (
	"context"
	"net"
	"sync"

	"github.com/webws/go-moda/logger"
	"github.com/webws/go-moda/transport"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
)

// 默认grpc address
const address = ":8082"

// 定义 ServerOptions func
type ServerOptions func(*Server)

func WithServerNetwork(network string) ServerOptions {
	return func(s *Server) {
		s.network = network
	}
}

func WithServerAddress(address string) ServerOptions {
	return func(s *Server) {
		s.address = address
	}
}

func WithTracing(tracing bool) ServerOptions {
	return func(s *Server) {
		s.tracing = tracing
	}
}

// 验证 Server 是否实现了 transport.Server 接口
var _ transport.Server = &Server{}

type Server struct {
	*grpc.Server
	ctx      context.Context
	listener net.Listener
	once     sync.Once
	network  string
	address  string
	tracing  bool
}

// NewServer
func NewServer(opts ...ServerOptions) *Server {
	srv := &Server{
		network: "tcp",
	}
	for _, o := range opts {
		o(srv)
	}
	if srv.address == "" {
		logger.Infow("[GRPC] server address is empty, use default address", "address", address)
		srv.address = address
	}
	var grpcOption []grpc.ServerOption

	srv.Server = grpc.NewServer()
	if srv.tracing {
		grpcOption = append(grpcOption, grpc.UnaryInterceptor(otelgrpc.UnaryServerInterceptor()))
		grpcOption = append(grpcOption, grpc.StreamInterceptor(otelgrpc.StreamServerInterceptor()))
	}
	srv.Server = grpc.NewServer(grpcOption...)
	return srv
}

// start
func (s *Server) Start(ctx context.Context) error {
	s.ctx = ctx
	var err error
	s.listener, err = net.Listen(s.network, s.address)
	if err != nil {
		logger.Errorw("[GRPC] server start failed", "listen_addr", s.address, "err", err)
		return err
	}
	logger.Infow("[GRPC] server started", "listen_addr", s.address)
	return s.Serve(s.listener)
}

// Stop
func (s *Server) Stop(ctx context.Context) error {
	s.once.Do(func() {
		s.Server.GracefulStop()
	})
	return nil
}
