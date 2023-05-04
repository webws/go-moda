package main

import (
	"net/http"

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
	pflag.StringVarP(&ConfFilePath, "conf", "c", "", "config file path")
	pflag.Parse()
	logger.Infow("helloworld", "conf", ConfFilePath, "server_name", ServerName, "app_version", AppVersion)
	logger.SetLevel(logger.InfoLevel)
	conf := &Config{}
	c := config.New(config.WithSources([]config.Source{
		&config.SourceFile{
			ConfigPath: ConfFilePath,
			// DefaultConfigPath: "/home/hellotalk/project/go_pro/go-note/go_moda/example/helloworld/conf.toml",
			DefaultConfigPath: "./conf.toml",
		},
		// &config.SourceText{"a=b"},
	}))
	if err := c.Load(conf); err != nil {
		panic(err)
	}
	hs := modahttp.NewNetHTTPServer()
	registerHttp(hs.GetServer())

	httpSrv := modahttp.NewServer(
		modahttp.WithAddress(conf.HttpAddr),
		modahttp.WitchHandle(hs),
	)
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
