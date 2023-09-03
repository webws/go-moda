package main

import (
	"context"
	"fmt"

	"github.com/spf13/pflag"
	"github.com/webws/go-moda/config"
	"github.com/webws/go-moda/logger"
	modagrpc "github.com/webws/go-moda/transport/grpc"

	app "github.com/webws/go-moda"
	pbexample "github.com/webws/go-moda/example/pb/example"
)

var (
	ServerName   string
	AppVersion   string
	ConfFilePath string
)

type Config struct {
	HttpAddr string `json:"http_addr" toml:"http_addr"`
	GrpcAddr string `json:"grpc_addr" toml:"grpc_addr"`
}

var csFlag = pflag.StringP("cs", "s", "client", "client or server")

func main() {
	// load config
	conf := &Config{}
	if err := config.NewConfigWithFile("./conf.toml").Load(conf); err != nil {
		logger.Fatalw("NewConfigWithFile fail", "err", err)
	}
	// 通过csFlag判断是启动服务端还是客户端
	if *csFlag == "server" {
		gs := modagrpc.NewServer(
			modagrpc.WithServerAddress(conf.GrpcAddr),
		)
		pbexample.RegisterExampleServiceServer(gs, &ExampleServer{})
		a := app.New(app.Server(gs))
		if err := a.Run(); err != nil {
			panic(err)
		}
	} else {
		ClientGrpcSend(conf.GrpcAddr)
	}
}

func ClientGrpcSend(addr string) {
	// 连接服务器
	conn, err := modagrpc.Dial(context.Background(), addr, modagrpc.WithDialWithInsecure(true), modagrpc.WithDialTracing(true))
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	// 创建一个grpc客户端
	client := pbexample.NewExampleServiceClient(conn)
	// 调用服务端的 SayHello 方法
	resp, err := client.SayHello(context.Background(), &pbexample.HelloRequest{Name: "gRPC"})
	if err != nil {
		panic(err)
	}
	fmt.Println(resp.Message)
}
