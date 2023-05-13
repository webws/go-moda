#### conf.toml
```toml
http_addr = ":8081"
grpc_addr = ":8082"
```
#### 运行
```shell
go run ./ -c ./conf.toml
```
* http 服务 http://localhost:8081/helloworld  
* grpc 服务 使用 gRPC 客户端调用 SayHello 方法