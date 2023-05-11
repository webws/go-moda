package main

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/spf13/pflag"
	"github.com/webws/go-moda/config"
	"github.com/webws/go-moda/logger"
	modahttp "github.com/webws/go-moda/transport/http"

	app "github.com/webws/go-moda"
	"github.com/webws/go-moda/tracing"
	// logger
)

type Config struct {
	HttpAddr  string `json:"http_addr" toml:"http_addr"`
	GrpcAddr  string `json:"grpc_addr" toml:"grpc_addr"`
	JaegerUrl string `json:"jaeger_url" toml:"jaeger_url"`
	Tracing   bool   `toml:"tracing"  json:"tracing"` // opentelemetry tracing
}

var (
	ServerName   = "api2"
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
		logger.Debugw("/api2/bar")
		ctx, span := tracing.Start(c.Request.Context(), "api2 bar_handler")
		defer span.End()
		_, err := tracing.CallAPI(ctx, "http://localhost:8083/api3/bar", "GET", nil)
		if err != nil {
			logger.Errorw("call api1 error", "err", err)
		}
		// Bar(c.Request.Context())
		c.JSON(http.StatusOK, http.StatusText(http.StatusOK))
	})
}
