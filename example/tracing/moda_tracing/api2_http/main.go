package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/webws/go-moda/config"
	"github.com/webws/go-moda/logger"
	modahttp "github.com/webws/go-moda/transport/http"

	app "github.com/webws/go-moda"
	configExample "github.com/webws/go-moda/example/config"
	"github.com/webws/go-moda/tracing"
	// logger
)

var (
	ServerName   = "api2"
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
	gin, httpSrv := modahttp.NewGinHttpServer(
		modahttp.WithAddress(conf.HttpAddr),
		modahttp.WithTracing(conf.Tracing),
	)
	registerHttp(gin)
	a := app.New(app.Server(httpSrv))
	a.Run()
}

func registerHttp(g *gin.Engine) {
	g.GET("/api2/bar", func(c *gin.Context) {
		// header
		log := tracing.LoggerWitchSpan(c.Request.Context(), logger.GetLogger())
		log.Infow("header", "header", c.Request.Header)
		ctx, span := tracing.Start(c.Request.Context(), "api2 bar_handler")
		defer span.End()
		log = tracing.LoggerWitchSpan(ctx, log)
		log.Infow("header", "header", c.Request.Header)
		url := fmt.Sprintf("http://%s/api3/bar", conf.ServiceAddr.Api3)
		_, err := modahttp.CallAPI(ctx, url, "GET", nil)
		if err != nil {
			logger.Errorw("call api1 error", "err", err)
		}
		// Bar(c.Request.Context())
		c.JSON(http.StatusOK, http.StatusText(http.StatusOK))
	})
}
