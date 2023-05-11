package main

import (
	"context"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
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
	ServerName   = "api3"
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
	e.GET("/api3/bar", func(c echo.Context) error {
		logger.Infow("/api1/bar info")
		time.Sleep(5 * time.Second)
		Bar(c.Request().Context())
		return c.JSON(http.StatusOK, http.StatusText(http.StatusOK))
	})
}
