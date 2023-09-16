package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/webws/go-moda/config"
	"github.com/webws/go-moda/logger"
	modahttp "github.com/webws/go-moda/transport/http"

	app "github.com/webws/go-moda"
	"github.com/webws/go-moda/tracing"

	modagrpc "github.com/webws/go-moda/transport/grpc"

	configExample "github.com/webws/go-moda/example/config"
	pbexample "github.com/webws/go-moda/example/pb/example"

	// logger
	"go.opentelemetry.io/otel/trace"
)

var (
	ServerName   = "api1"
	AppVersion   string
	ConfFilePath string
)
var conf *configExample.Config

func main() {
	conf = &configExample.Config{}

	if err := config.NewConfigWithFile("./conf.toml").Load(conf); err != nil {
		logger.Fatalw("NewConfigWithFile fail", "err", err)
	}
	conf.SetEnvServiceAddr()
	// init jaeger provider
	shutdown, err := tracing.InitJaegerProvider(conf.JaegerUrl, ServerName)
	if err != nil {
		panic(err)
	}
	defer shutdown(context.Background())
	e, httpSrv := modahttp.NewEchoHttpServer(
		modahttp.WithAddress(conf.HttpAddr),
	)
	registerHttp(e)
	a := app.New(app.Server(httpSrv))
	a.Run()
}

func registerHttp(e *echo.Echo) {
	e.GET("/api1/bar", func(c echo.Context) error {
		logger.Infow("/api1/bar info", "req.header", c.Request().Header)
		// l:=logger.With("api1", "abc")

		spanCtx := trace.SpanContextFromContext(c.Request().Context())
		logger.Infow("span", "span_id", spanCtx.SpanID().String(), "trace_id", spanCtx.TraceID().String())
		// call api2
		ctx, span := tracing.Start(c.Request().Context(), "api1")
		defer span.End()

		url := fmt.Sprintf("http://%s/api2/bar", conf.ServiceAddr.Api2)
		spanCtx = trace.SpanContextFromContext(ctx)
		logger.Infow("span", "span_id", spanCtx.SpanID().String(), "trace_id", spanCtx.TraceID().String())

		defer span.End()
		_, err := modahttp.CallAPI(ctx, url, "GET", nil)
		if err != nil {
			logger.Errorw("call api2 error", "err", err)
		}
		ClientGrpcSend(conf.ServiceAddr.Grpc1, c.Request().Context())
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
