package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/spf13/pflag"
	app "github.com/webws/go-moda"
	"github.com/webws/go-moda/config"
	"github.com/webws/go-moda/logger"
	modahttp "github.com/webws/go-moda/transport/http"
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
	pflag.StringVarP(&ConfFilePath, "conf", "c", "", "config file path")
	pflag.Parse()
	logger.Infow("helloworld", "conf", ConfFilePath, "server_name", ServerName, "app_version", AppVersion)
	logger.SetLevel(logger.InfoLevel)
	conf := &Config{}
	c := config.New(config.WithSources([]config.Source{
		&config.SourceFile{
			ConfigPath: ConfFilePath,
			DefaultConfigPath: "./conf.toml",
		},
		// &config.SourceText{"a=b"},
	}))
	if err := c.Load(conf); err != nil {
		panic(err)
	}
	// gin http
	hs := modahttp.NewGinServer()
	registerHttp(hs.GetServer())
	httpSrv := modahttp.NewServer(
		modahttp.WithAddress(conf.HttpAddr),
		modahttp.WitchHandle(hs),
	)
	// app run
	a := app.New(app.Server(httpSrv))
	a.Run()
}

func registerHttp(g *gin.Engine) {
	g.GET("/helloworld", func(c *gin.Context) {
		logger.Debugw("helloworld debug")
		c.JSON(http.StatusOK, http.StatusText(http.StatusOK))
	})
}
