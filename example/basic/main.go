package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	app "github.com/webws/go-moda"
	"github.com/webws/go-moda/config"
	examplepb "github.com/webws/go-moda/example/pb/example"
	"github.com/webws/go-moda/logger"
	"google.golang.org/grpc"
)

var ServerName string

type Config struct {
	HttpAddr string `json:"http_addr" toml:"http_addr"`
	GrpcAddr string `json:"grpc_addr" toml:"grpc_addr"`
}

func main() {
	logger.SetLevel(logger.DebugLevel)
	conf := &Config{}
	if err := config.NewConfigWithFile("./conf.toml").Load(conf); err != nil {
		logger.Fatalw("NewConfigWithFile fail", "err", err)
	}
	// http server default gin
	a := app.NewServer().AddHttpServer(conf.HttpAddr, registerHttp)
	// grpc
	a.AddGrpcServer(conf.GrpcAddr, func(sr grpc.ServiceRegistrar) {
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
