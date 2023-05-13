# go-moda
go-moda 是一个基于 Go 语言的通用开发框架
## 特性
- 统一启动入口和优雅退出
- config:    通用的配置文件读取模块，支持 toml、yaml 和 json 格式。
- logger:    日志系统模块，基于 Zap,并支持全局日志和模块日志。
- pprof:	 分析程序的工具
- transport: 集成 http（echo、gin）和 grpc。
- tracing:   openTelemetry 实现微务链路追踪(TODO)
- Metrics:   指标系统,集成 Prometheus (TODO)

## 快速使用
### conf.toml
```toml
http_addr = ":8081"
grpc_addr = ":8082"
```
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
	HttpAddr  string `json:"http_addr" toml:"http_addr"`
	GrpcAddr  string `json:"grpc_addr" toml:"grpc_addr"`
	JaegerUrl string `json:"jaeger_url" toml:"jaeger_url"`
	Tracing   bool   `toml:"tracing"  json:"tracing"` // opentelemetry tracing
}

var ServerName = "example"

// AppVersion   string
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
		logger.Fatalw("load config error", "err", err)
	}
	// http server
	gin, httpSrv := modahttp.NewGinHttpServer(
		modahttp.WithAddress(conf.HttpAddr),
	)
	registerHttp(gin)

	// grpc server
	grpcSrv := modagrpc.NewServer(modagrpc.WithServerAddress(conf.GrpcAddr))
	grecExample := &ExampleServer{}
	pbexample.RegisterExampleServiceServer(grpcSrv, grecExample)

	// app run
	a := app.New(
		app.Server(httpSrv, grpcSrv),
		app.Name(ServerName),
	)
	if err := a.Run(); err != nil {
		logger.Fatalw("app run error", "err", err)
	}
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

func (s *ExampleServer) SayHello(ctx context.Context, req *pbexample.HelloRequest) (*pbexample.HelloResponse, error) {
	return &pbexample.HelloResponse{Message: "Hello " + req.Name}, nil
}

```
### 运行
```shell
go run ./ -c ./conf.toml
```
* http 服务 http://localhost:8081/helloworld  
* grpc 服务 使用 gRPC 客户端调用 SayHello 方法 
## pprof 性能分析
启动服务默认开启 pprof 性能分析，可以通过 http://localhost:8081/debug/pprof/ 查看
![](images/2023-05-13-11-02-02.png)
可视化分析 gouroutine
```shell
go tool pprof http://localhost:8081/debug/pprof/goroutine
(pprof) web
```
可能提示 需要先安装 graphviz, mac 下可以使用 brew 安装
```shell
brew install graphviz
```
![](images/2023-05-13-11-04-41.png)

## tracing
* moda-go 集成了 opentelemetry 实现微服务链路追踪，目前 exporter 支持 jaeger 
* 演示启动的是 jaeger all-in-one: 详细示例请看:[tracing_example](./example/tracing/moda_tracing/)

![](./example/tracing/moda_tracing/images/2023-05-12-01-08-57.png)