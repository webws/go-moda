package main

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	app "github.com/webws/go-moda"
	pbexample "github.com/webws/go-moda/example/pb/example"
	"github.com/webws/go-moda/logger"
	"github.com/webws/go-moda/tracing"
	"google.golang.org/grpc"
)

var ServerName string

func main() {
	// http server default gin
	a := app.NewServer().AddHttpServer(":8088", registerHttp)
	// grpc
	a.AddGrpcServer(":8087", func(sr grpc.ServiceRegistrar) {
		pbexample.RegisterExampleServiceServer(sr, &ExampleServer{})
	})
	// tracing
	// Need to start jaeger
	shutdown, err := tracing.InitJaegerProvider("http://localhost:14268/api/traces", ServerName)
	if err != nil {
		logger.Fatalw("InitJaegerProvider", "err", err)
	}
	defer shutdown(context.Background())
	a.SetTracing(true)
	if err := a.Run(); err != nil {
		logger.Fatalw("app run error", "err", err)
	}
}

func registerHttp(g *gin.Engine) {
	g.GET("/helloworld", func(c *gin.Context) {
		logger.Debugw("Hello World")
		c.JSON(http.StatusOK, http.StatusText(http.StatusOK))
	})
}

type ExampleServer struct {
	pbexample.UnimplementedExampleServiceServer
}

func (s *ExampleServer) SayHello(ctx context.Context, req *pbexample.HelloRequest) (*pbexample.HelloResponse, error) {
	return &pbexample.HelloResponse{Message: "Hello " + req.Name}, nil
}
