

## 服务链路关系
#### 关系图
<!-- 调用架构图 -->
```mermaid
graph LR
A[用户] --> B[api1/bar]
B --> C[api2/bar]
C --> D[api3/bar]
D --> E[bar]
E --> F[bar2]
F --> G[bar3]
```

#### 关系说明:
1. 用户  请求 api1(echo server) 服务的 api1/bar
2. api1 调用 api2 (gin server) 服务的 api2/bar
3. api2 调用 api3 (echo server )服务的 api3/bar
4. api3 调用 内部 调用方法 bar->bar2->bar3

## 安装jaeger
1. 下载jaeger,我使用的是 jaeger-all-in-one
2. 启动 jaeger ~/tool/jaeger-1.31.0-linux-amd64/jaeger-all-in-one
3. 默认查看面板 地址 http://localhost:16686/
4. tracer Batcher的地址,下面代码会体现: http://localhost:14268/api/traces

## 调用服务,查看链路关系
### 运行代码示例
1. 示例文件:moda_tracing下 有三个目录,分别是 api1_http,api2_http,api_http,分别对应三个服务  
2. 分别启动三个服务,进入目录 go run ./ 即可启动服务,端口分别是 8081,8082,8083
3. 根据上面链路关系,调用api1: curl localhost:8081/api1/bar
4. 等待调用成功
5. 打开 jaeger 面板,查看链路关系图,http://localhost:16686/
6. todo:示例代码启动采用 docker-compose 启动,方便演示
   
### 查看jaeger 
![](images/2023-05-12-01-08-57.png)





