#go-moda
 golang 通用的 grpc http 基础开发框架
## 特性
- transport: 集成 http（echo、gin）和 grpc。
- 统一服务启动入口
- config:    通用的配置文件读取模块，支持 toml、yaml 和 json 格式。
- logger:    结构化 统一 logger API, 已 新增 slog 替换zap [logger](./logger/)
- pprof:	 分析性能
- tracing:   openTelemetry 实现微务链路追踪
### 快速使用

``` golang
package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	app "github.com/webws/go-moda"
	examplepb "github.com/webws/go-moda/example/pb/example"

	"github.com/webws/go-moda/logger"
	"google.golang.org/grpc"
)

var ServerName string

func main() {
	// http server default gin
	a := app.NewServer().AddHttpServer(":8081", registerHttp)
	// grpc
	a.AddGrpcServer(":8082", func(sr grpc.ServiceRegistrar) {
		examplepb.RegisterExampleServiceServer(sr, &examplepb.ExampleServer{})
	})
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

```

#### 运行
```shell
go run ./ 
```
* 请求 http url http://localhost:8081/helloworld  
* grpc 服务 使用 gRPC 客户端调用 SayHello 方法 [example_grpc](./example/grpc/)

### config

``` golang
// "github.com/webws/go-moda/config"
//...
	conf := &Config{}
	if err := config.NewConfigWithFile("./conf.toml").Load(conf); err != nil {
		logger.Fatalw("NewConfigWithFile fail", "err", err)
	}
//...
```
### logger
``` golang
package main

import "github.com/webws/go-moda/logger"

func main() {
	// 格式化打印 {"time":"2023-09-08T01:25:21.313463+08:00","level":"INFO","msg":"info hello slog","key":"value","file":"/Users/xxx/w/pro/go-moda/example/logger/main.go","line":6}
	logger.Infow("info hello slog", "key", "value")   // 打印json
	logger.Debugw("debug hello slog", "key", "value") // 不展示
	logger.SetLevel(logger.DebugLevel)                // 设置等级
	logger.Debugw("debug hello slog", "key", "value") // 设置了等级之后展示 debug
	// with
	newLog := logger.With("newkey", "newValue")
	newLog.Debugw("new hello slog") // 会打印 newkey:newValue
	logger.Debugw("old hello slog") // 不会打印 newkey:newValue
}
``` 


## tracing 链路追踪
* 使用 opentelemetry 实现微服务链路追踪，目前 exporter 支持 jaeger 
* 示例集成了docker 环境,支持 make deploy 同时启动 jaeger,api1,api2,api3,grpc 服务
* 详细示例请看:[tracing_example](./example/tracing/moda_tracing/)

1. server 端使用
```golang
var ServerName string
func main() {
	// http server default gin
	a := app.NewServer().AddHttpServer(":8081", registerHttp)
	// grpc
	a.AddGrpcServer(":8082", func(sr grpc.ServiceRegistrar) {
		pbexample.RegisterExampleServiceServer(sr, &ExampleServer{})
	})
	// tracing
	// Need to start jaeger
	shutdown, err := tracing.InitJaegerProvider("http://localhost:14268/api/traces", ServerName)
	if err != nil {
		logger.Fatalw("InitJaegerProvider", "err", err)
	}
	defer shutdown(context.Background())
	a.SetTracing(true)
	if err := a.Run(); err != nil {
		logger.Fatalw("app run error", "err", err)
	}
}
```  
* tracing.InitJaegerProvider 初始化 Jaeger Provider
* a.SetTracing(true) 设置tracing为true,内部将 请求头的参数  转换成 ctx

2.  client端:  封装了 CallAPI 方法, 已将span ctx 信息注入到请求头
```golang
    // ...
    _, err := modahttp.CallAPI(ctx, url, "POST", nil)
		
```

![](./example/tracing/moda_tracing/images/2023-05-13-23-02-24.png)
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