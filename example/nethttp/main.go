package main

import (
	"net/http"

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
	// load config
	conf := &Config{}
	if err := config.NewConfigWithFile("./conf.toml").Load(conf); err != nil {
		logger.Fatalw("NewConfigWithFile fail", "err", err)
	}
	serveMux, httpSrv := modahttp.NewNetHttpServer(modahttp.WithAddress(conf.HttpAddr))
	registerHttp(serveMux)
	a := app.New(app.Server(httpSrv))
	a.Run()
}

func registerHttp(serve *http.ServeMux) {
	// hello world
	serve.HandleFunc("/helloworld", func(w http.ResponseWriter, r *http.Request) {
		logger.Debugw("helloworld debug")
		w.Write([]byte(http.StatusText(http.StatusOK)))
	})
}
