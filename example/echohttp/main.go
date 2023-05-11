package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/spf13/pflag"
	"github.com/webws/go-moda/config"
	"github.com/webws/go-moda/logger"
	modahttp "github.com/webws/go-moda/transport/http"

	app "github.com/webws/go-moda"
)

type Config struct {
	HttpAddr string `json:"http_addr" toml:"http_addr"`
	GrpcAddr string `json:"grpc_addr" toml:"grpc_addr"`
}

var (
	ServerName   string
	AppVersion   string
	ConfFilePath string
)

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
		// &config.SourceText{"a=b"},
	}))
	if err := c.Load(conf); err != nil {
		panic(err)
	}
	// http server
	e, httpSrv := modahttp.NewEchoHttpServer(modahttp.WithAddress(conf.HttpAddr))
	registerHttp(e)

	// run
	a := app.New(app.Server(httpSrv))
	if err := a.Run(); err != nil {
		panic(err)
	}
}

func registerHttp(e *echo.Echo) {
	e.GET("/helloworld", func(c echo.Context) error {
		logger.Debugw("helloworld debug")
		return c.JSON(http.StatusOK, http.StatusText(http.StatusOK))
	})
}
