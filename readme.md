# go-moda

go-moda 是一个基于 Go 语言开发的通用开发框架
## 特性
- Config：通用的配置文件读取模块，支持 toml、yaml 和 json 格式。
- Logger：日志系统模块，基于 Zap，并支持全局日志和模块日志。
- Transport：HTTP（Echo、Gin 和 net/http）和 GRPC。
- 统一启动入口和优雅退出
- Pprof
- sentry (待实现)
- Prometheus (待实现)
- Tracing (http server and client 已封装) (grpc 待实现)
- Makefile local build,ci/cd前置处理 (待实现)

## 快速使用
### main.go
```golang
package main

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/spf13/pflag"
	app "github.com/webws/go-moda"
	"github.com/webws/go-moda/config"
	pbexample "github.com/webws/go-moda/example/pb/example"
	"github.com/webws/go-moda/logger"
	modagrpc "github.com/webws/go-moda/transport/grpc"
	modahttp "github.com/webws/go-moda/transport/http"
)

type Config struct {
	HttpAddr string `json:"http_addr" toml:"http_addr"`
	GrpcAddr string `json:"grpc_addr" toml:"grpc_addr"`
}

var ConfFilePath string

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
	}))
	if err := c.Load(conf); err != nil {
		panic(err)
	}
	// http server
	gin, httpSrv := modahttp.NewGinHttpServer(
		modahttp.WithAddress(conf.HttpAddr),
	)
	registerHttp(gin)

	// grpc server
	grpcSrv := modagrpc.NewServer()
	grecExample := &ExampleServer{}
	pbexample.RegisterExampleServiceServer(grpcSrv, grecExample)
	// app run
	a := app.New(app.Server(httpSrv, grpcSrv))
	a.Run()
}

func registerHttp(g *gin.Engine) {
	g.GET("/helloworld", func(c *gin.Context) {
		logger.Debugw("helloworld debug")
		c.JSON(http.StatusOK, http.StatusText(http.StatusOK))
	})
}

// 实现 GrpcServer 接口
type ExampleServer struct {
	pbexample.UnimplementedExampleServiceServer
}

// 实现 GrpcServer 接口中的 SayHello 方法

func (s *ExampleServer) SayHello(ctx context.Context, req *pbexample.HelloRequest) (*pbexample.HelloResponse, error) {
	// return nil, nil
	return &pbexample.HelloResponse{Message: "Hello " + req.Name}, nil
}


```
### conf.toml
```toml
http_addr = ":8081"
grpc_addr = ":8082"
```
### 运行
```shell
go run ./ -c ./conf.toml
```
### 服务启动后
1. http 服务 http://localhost:8081/helloworld
2. grpc 服务 使用 gRPC 客户端调用 SayHello 方法
3. pprof http://localhost:8081/debug/pprof/
## tracing
tracing 示例: [tracing](./example/tracing/moda_tracing/)