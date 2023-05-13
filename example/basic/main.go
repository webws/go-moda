package main

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/spf13/pflag"
	app "github.com/webws/go-moda"
	"github.com/webws/go-moda/config"
	pbexample "github.com/webws/go-moda/example/pb/example"
	"github.com/webws/go-moda/logger"
	modagrpc "github.com/webws/go-moda/transport/grpc"
	modahttp "github.com/webws/go-moda/transport/http"
)

type Config struct {
	HttpAddr  string `json:"http_addr" toml:"http_addr"`
	GrpcAddr  string `json:"grpc_addr" toml:"grpc_addr"`
	JaegerUrl string `json:"jaeger_url" toml:"jaeger_url"`
	Tracing   bool   `toml:"tracing"  json:"tracing"` // opentelemetry tracing
}

var ServerName = "example"

// AppVersion   string
var ConfFilePath string

func main() {
	// flag
	pflag.StringVarP(&ConfFilePath, "conf", "c", "", "config file path")
	pflag.Parse()
	// set logger level info,default is debug
	logger.SetLevel(logger.InfoLevel)
	// load config
	conf := &Config{}
	c := config.New(config.WithSources([]config.Source{
		&config.SourceFile{
			ConfigPath:        ConfFilePath,
			DefaultConfigPath: "./conf.toml",
		},
	}))
	if err := c.Load(conf); err != nil {
		logger.Fatalw("load config error", "err", err)
	}
	// http server
	gin, httpSrv := modahttp.NewGinHttpServer(
		modahttp.WithAddress(conf.HttpAddr),
	)
	registerHttp(gin)

	// grpc server
	grpcSrv := modagrpc.NewServer(modagrpc.WithServerAddress(conf.GrpcAddr))
	grecExample := &ExampleServer{}
	pbexample.RegisterExampleServiceServer(grpcSrv, grecExample)

	// app run
	a := app.New(
		app.Server(httpSrv, grpcSrv),
		app.Name(ServerName),
	)
	if err := a.Run(); err != nil {
		logger.Fatalw("app run error", "err", err)
	}
}

func registerHttp(g *gin.Engine) {
	g.GET("/helloworld", func(c *gin.Context) {
		logger.Debugw("helloworld debug")
		c.JSON(http.StatusOK, http.StatusText(http.StatusOK))
	})
}

// 实现 GrpcServer 接口
type ExampleServer struct {
	pbexample.UnimplementedExampleServiceServer
}

func (s *ExampleServer) SayHello(ctx context.Context, req *pbexample.HelloRequest) (*pbexample.HelloResponse, error) {
	return &pbexample.HelloResponse{Message: "Hello " + req.Name}, nil
}
