package main

import (
	"context"
	"net"

	"github.com/spf13/pflag"
	"github.com/webws/go-moda/config"
	"github.com/webws/go-moda/logger"
	"github.com/webws/go-moda/tracing"
	modagrpc "github.com/webws/go-moda/transport/grpc"
	"google.golang.org/grpc"

	app "github.com/webws/go-moda"
	pbexample "github.com/webws/go-moda/example/pb/example"

	configExample "github.com/webws/go-moda/example/config"
)

var (
	ServerName   string
	AppVersion   string
	ConfFilePath string
)

var csFlag = pflag.StringP("cs", "s", "client", "client or server")

var conf *configExample.Config

func main() {
	conf = &configExample.Config{}

	if err := config.NewConfigWithFile("./conf.toml").Load(conf); err != nil {
		logger.Fatalw("NewConfigWithFile fail", "err", err)
	}
	conf.SetEnvServiceAddr()
	// init jaeger provider
	shutdown, err := tracing.InitJaegerProvider(conf.JaegerUrl, "grpc-server")
	if err != nil {
		panic(err)
	}
	defer shutdown(context.Background())
	gs := modagrpc.NewServer(
		modagrpc.WithServerAddress(conf.GrpcAddr),
		modagrpc.WithTracing(conf.Tracing),
	)
	pbexample.RegisterExampleServiceServer(gs, &ExampleServer{})
	a := app.New(app.Server(gs))
	if err := a.Run(); err != nil {
		panic(err)
	}
}

// 启动一个grpc server
func StartGrpcServer() {
	// listen
	lis, err := net.Listen("tcp", ":8082")
	if err != nil {
		panic(err)
	}
	// 创建一个grpc server
	gs := grpc.NewServer()

	// 注册服务
	pbexample.RegisterExampleServiceServer(gs, &ExampleServer{})
	// start
	logger.Infow("[GRPC] server started", "listen_addr", ":8082")
	if err := gs.Serve(lis); err != nil {
		panic(err)
	}
}
