package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	app "github.com/webws/go-moda"
	examplepb "github.com/webws/go-moda/example/pb/example"

	"github.com/webws/go-moda/logger"
	"google.golang.org/grpc"
)

var ServerName string

func main() {
	// http server default gin
	a := app.NewServer().AddHttpServer(":8088", registerHttp)
	// grpc
	a.AddGrpcServer(":8087", func(sr grpc.ServiceRegistrar) {
		examplepb.RegisterExampleServiceServer(sr, &examplepb.ExampleServer{})
	})
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
