package main

import (
	"context"
	"fmt"
	"net"

	"github.com/spf13/pflag"
	"github.com/webws/go-moda/config"
	"github.com/webws/go-moda/logger"
	modagrpc "github.com/webws/go-moda/transport/grpc"
	"google.golang.org/grpc"

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
	pflag.StringVarP(&ConfFilePath, "conf", "c", "", "config file path")
	pflag.Parse()
	conf := &Config{}
	c := config.New(config.WithSources([]config.Source{
		&config.SourceFile{
			ConfigPath:        ConfFilePath,
			DefaultConfigPath: "./conf.toml",
		},
		// &config.SourceText{"a=b"},
	}))
	if err := c.Load(conf); err != nil {
		panic(err)
	}
	// 通过csFlag判断是启动服务端还是客户端
	if *csFlag == "server" {
		gs := modagrpc.NewServer(
			modagrpc.WithServerAddress(conf.GrpcAddr),
			modagrpc.WithServerNetwork("tcp"),
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
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
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
