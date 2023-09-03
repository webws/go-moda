package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	app "github.com/webws/go-moda"
	"github.com/webws/go-moda/config"
	"github.com/webws/go-moda/logger"
	modahttp "github.com/webws/go-moda/transport/http"
)

type Config struct {
	HttpAddr string `json:"http_addr" toml:"http_addr"`
	GrpcAddr string `json:"grpc_addr" toml:"grpc_addr"`
}

func main() {
	conf := &Config{}
	if err := config.NewConfigWithFile("./conf.toml").Load(conf); err != nil {
		logger.Fatalw("NewConfigWithFile fail", "err", err)
	}
	// gin http server
	gin, httpSrv := modahttp.NewGinHttpServer(modahttp.WithAddress(conf.HttpAddr))
	registerHttp(gin)

	// app run
	a := app.New(app.Server(httpSrv))
	if err := a.Run(); err != nil {
		logger.Fatalw("app run fail", "err", err)
	}
}

func registerHttp(g *gin.Engine) {
	g.GET("/", func(c *gin.Context) {
		logger.Debugw("helloworld debug")
		c.JSON(http.StatusOK, http.StatusText(http.StatusOK))
	})
}
