package grpc

import (
	"context"
	"net"
	"net/url"
	"sync"

	"github.com/webws/go-moda/logger"
	"github.com/webws/go-moda/transport"
	"google.golang.org/grpc"
)

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

// 验证 Server 是否实现了 transport.Server 接口
var _ transport.Server = &Server{}

type Server struct {
	*grpc.Server
	ctx      context.Context
	listener net.Listener
	once     sync.Once
	network  string
	address  string
	endpoint *url.URL
}

// NewServer
func NewServer(opts ...ServerOptions) *Server {
	srv := &Server{
		Server:  grpc.NewServer(),
		network: "tcp",
		address: ":8082",
	}
	for _, o := range opts {
		o(srv)
	}
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

// endpoint
func (s *Server) Endpoint() (*url.URL, error) {
	// TODO Endpoint用于非k8s环境下的服务发现
	s.endpoint = &url.URL{
		Scheme: "grpc",
		Host:   s.listener.Addr().String(),
	}
	return s.endpoint, nil
}
