package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/webws/go-moda/config"
	"github.com/webws/go-moda/logger"
	modahttp "github.com/webws/go-moda/transport/http"

	app "github.com/webws/go-moda"
)

type Config struct {
	HttpAddr string `json:"http_addr" toml:"http_addr"`
}

func main() {
	// load config
	conf := &Config{}
	if err := config.NewConfigWithFile("./conf.toml").Load(conf); err != nil {
		logger.Fatalw("NewConfigWithFile fail", "err", err)
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
