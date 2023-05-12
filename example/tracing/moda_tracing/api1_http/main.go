package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/spf13/pflag"
	"github.com/webws/go-moda/config"
	"github.com/webws/go-moda/logger"
	modahttp "github.com/webws/go-moda/transport/http"

	app "github.com/webws/go-moda"
	"github.com/webws/go-moda/tracing"

	modagrpc "github.com/webws/go-moda/transport/grpc"

	pbexample "github.com/webws/go-moda/example/pb/example"
	// logger
)

type Config struct {
	HttpAddr  string `json:"http_addr" toml:"http_addr"`
	GrpcAddr  string `json:"grpc_addr" toml:"grpc_addr"`
	JaegerUrl string `json:"jaeger_url" toml:"jaeger_url"`
	Tracing   bool   `toml:"tracing"  json:"tracing"` // opentelemetry tracing
}

var (
	ServerName   = "api1"
	AppVersion   string
	ConfFilePath string
)

func main() {
	pflag.StringVarP(&ConfFilePath, "conf", "c", "", "config file path")
	pflag.Parse()
	// load config
	logger.SetLevel(logger.DebugLevel)
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
	// init jaeger provider
	shutdown, err := tracing.InitJaegerProvider(conf.JaegerUrl, ServerName)
	if err != nil {
		panic(err)
	}
	defer shutdown(context.Background())
	e, httpSrv := modahttp.NewEchoHttpServer(
		modahttp.WithAddress(conf.HttpAddr),
		modahttp.WithTracing(conf.Tracing),
	)
	registerHttp(e)
	a := app.New(app.Server(httpSrv))
	a.Run()
}

func registerHttp(e *echo.Echo) {
	e.GET("/api1/bar", func(c echo.Context) error {
		logger.Infow("/api1/bar info")
		// call api2
		_, err := modahttp.CallAPI(c.Request().Context(), "http://localhost:8082/api2/bar", "GET", nil)
		if err != nil {
			logger.Errorw("call api2 error", "err", err)
		}
		ClientGrpcSend("localhost:8086", c.Request().Context())
		// Bar(c.Request().Context())
		return c.JSON(http.StatusOK, http.StatusText(http.StatusOK))
	})
}

func ClientGrpcSend(addr string, ctx context.Context) {
	// 连接服务器
	conn, err := modagrpc.Dial(ctx, addr, modagrpc.WithDialWithInsecure(true), modagrpc.WithDialTracing(true))
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	// 创建一个grpc客户端
	client := pbexample.NewExampleServiceClient(conn)
	// 调用服务端的 SayHello 方法
	resp, err := client.SayHello(ctx, &pbexample.HelloRequest{Name: "gRPC"})
	if err != nil {
		panic(err)
	}
	fmt.Println(resp.Message)
}
