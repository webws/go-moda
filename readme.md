# go-moda
 golang 通用的 grpc http 开发框架,持续更新中
## 特性
- transport: 集成 http（echo、gin）和 grpc。
- 统一服务启动入口
- config:    通用的配置文件读取模块，支持 toml、yaml 和 json 格式。
- logger:    结构化 统一 logger API, 已 新增 slog 替换zap [logger](./logger/)
- pprof:	 分析性能
- tracing:   openTelemetry 实现微务链路追踪
- Metrics:   指标系统,集成 Prometheus (TODO)
## 快速使用
#### conf.toml
```toml
http_addr = ":8081"
grpc_addr = ":8082"
```
#### 启用http(gin) 和 grpc服务

``` golang
package main

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	app "github.com/webws/go-moda"
	"github.com/webws/go-moda/config"
	pbexample "github.com/webws/go-moda/example/pb/example"
	"github.com/webws/go-moda/logger"
	modagrpc "github.com/webws/go-moda/transport/grpc"
	modahttp "github.com/webws/go-moda/transport/http"
)

var ServerName string

type Config struct {
	HttpAddr string `json:"http_addr" toml:"http_addr"`
	GrpcAddr string `json:"grpc_addr" toml:"grpc_addr"`
}

func main() {
	conf := &Config{}
	if err := config.NewConfigWithFile("./conf.toml").Load(conf); err != nil {
		logger.Fatalw("NewConfigWithFile fail", "err", err)
	}
	// http server
	gin, httpSrv := modahttp.NewGinHttpServer(
		modahttp.WithAddress(conf.HttpAddr),
	)
	registerHttp(gin)

	// grpc server
	grpcSrv := modagrpc.NewServer(
		modagrpc.WithServerAddress(conf.GrpcAddr),
	)
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
		logger.Debugw("Hello World")
		c.JSON(http.StatusOK, http.StatusText(http.StatusOK))
	})
}

type ExampleServer struct {
	pbexample.UnimplementedExampleServiceServer
}

func (s *ExampleServer) SayHello(ctx context.Context, req *pbexample.HelloRequest) (*pbexample.HelloResponse, error) {
	return &pbexample.HelloResponse{Message: "Hello " + req.Name}, nil
}

```

#### 运行
```shell
go run ./ -c ./conf.toml
``````
* 请求 http url http://localhost:8081/helloworld  
* grpc 服务 使用 gRPC 客户端调用 SayHello 方法

更多服务启用示例
1. echo http :[example_echo](./example/echohttp/)
2. net http :[example_echo](./example/nethttp/)
3. grpc [example_grpc](./example/grpc/)
## pprof 性能分析
启动服务默认开启 pprof 性能分析，浏览器打开 http://localhost:8081/debug/ 查看
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

## tracing 链路追踪
* 使用 opentelemetry 实现微服务链路追踪，目前 exporter 支持 jaeger 
* 示例集成了docker 环境,支持 make deploy 同时启动 jaeger,api1,api2,api3,grpc 服务
* 详细示例请看:[tracing_example](./example/tracing/moda_tracing/)

1. 初始化 jaeger tracing  
```
import "github.com/webws/go-moda/tracing"
func main(){
    //...
    shutdown, err := tracing.InitJaegerProvider(conf.JaegerUrl, "grpc-server")
	if err != nil {
		panic(err)
	}
	defer shutdown(context.Background())
    //...
}
```
2. 在代码主动tracing start
```
  ctx, span := tracing.Start(c.Request().Context(), "api1")
  defer span.End()
```
3. 服务之间调用 产生的链路
   
*  server端: 增加 WithTracing 即可
```
    //...
    gin, httpSrv := modahttp.NewGinHttpServer(
		modahttp.WithAddress(conf.HttpAddr),
		modahttp.WithTracing(true),
	)
```
 * client端:  封装了 CallAPI 方法, 已将span ctx 信息注入到请求头
```
    // ...
    _, err := modahttp.CallAPI(ctx, url, "POST", nil)
		
```


![](./example/tracing/moda_tracing/images/2023-05-13-23-02-24.png)
## 更多示例
* 基本http/grpc服务启动示例:[basic example](./example/basic/)
* gin http 服务示例:[example_gin](./example/ginhttp/)
* echo http 服务示例:[example_echo](./example/echohttp/)
* grpc 服务示例:[example_grpc](./example/grpc/)
* tracing 示例:[tracing_example](./example/tracing/moda_tracing/)
* docker 本地CD示例:[docker_example](./example/tracing/moda_tracing/)
## 参考链接
* https://github.com/open-telemetry
* https://github.com/go-kratos
* https://github.com/open-telemetry/opentelemetry-go-contrib
* https://github.com/labstack/echo
* https://github.com/gin-gonic/gin
* https://github.com/uber-go/zap