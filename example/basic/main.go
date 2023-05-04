package main

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/spf13/pflag"
	app "github.com/webws/go-moda"
	"github.com/webws/go-moda/config"
	pbexample "github.com/webws/go-moda/example/pb/example"
	"github.com/webws/go-moda/logger"
	modagrpc "github.com/webws/go-moda/transport/grpc"
	modahttp "github.com/webws/go-moda/transport/http"
)

type Config struct {
	HttpAddr string `json:"http_addr" toml:"http_addr"`
	GrpcAddr string `json:"grpc_addr" toml:"grpc_addr"`
}

// ServerName   string
// AppVersion   string
var ConfFilePath string

func main() {
	// flag
	pflag.StringVarP(&ConfFilePath, "conf", "c", "", "config file path")
	pflag.Parse()
	// set logger level info,default is debug
	logger.Debugw("debug1", "debug", "debug")
	logger.SetLevel(logger.InfoLevel)
	logger.Debugw("debug2", "debug", "debug")
	// load config
	conf := &Config{}
	c := config.New(config.WithSources([]config.Source{
		&config.SourceFile{
			ConfigPath:        ConfFilePath,
			DefaultConfigPath: "./conf.toml",
		},
	}))
	if err := c.Load(conf); err != nil {
		panic(err)
	}
	// http server
	echoServer := modahttp.NewEchoServer()
	registerHttp(echoServer.GetServer())
	httpSrv := modahttp.NewServer(
		modahttp.WithAddress(conf.HttpAddr),
		modahttp.WitchHandle(echoServer),
	)
	// grpc server
	grpcSrv := modagrpc.NewServer()
	grecExample := &ExampleServer{}
	pbexample.RegisterExampleServiceServer(grpcSrv, grecExample)
	// app run
	a := app.New(app.Server(httpSrv, grpcSrv))
	a.Run()
}

func registerHttp(e *echo.Echo) {
	e.GET("/helloworld", func(c echo.Context) error {
		logger.Debugw("helloworld debug")
		return c.JSON(http.StatusOK, http.StatusText(http.StatusOK))
	})
}

// 实现 GrpcServer 接口
type ExampleServer struct {
	pbexample.UnimplementedExampleServiceServer
}

// 实现 GrpcServer 接口中的 SayHello 方法

func (s *ExampleServer) SayHello(ctx context.Context, req *pbexample.HelloRequest) (*pbexample.HelloResponse, error) {
	// return nil, nil
	return &pbexample.HelloResponse{Message: "Hello " + req.Name}, nil
}
