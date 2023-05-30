### 本机开发部署小工具
1. 本机开发多个服务,每个服务还会互相调用
2. 正常情况下,开发者需要在本地启动多个服务,并且手动调用
3. 不使用gitlab ci/cd, 本机通过 go build+makefile+docker-compose 编排多个服务
### 正常cicd流程
```mermaid
graph LR
A(开发者) -- push --> B[GitLab ci/cd]
B -- go build --> F[镜像地址更新]
F -- 触发部署 --> H[k8s/docker sync 运行新功能]
```
### 本机 makefile +docker 部署流程
1. 开发者写好功能,本地 go build 为二进制包
2. dockerfile 基于 alpha 镜像,运行二进制包
3. docker-compose 编排执行 dockerfile 运行多个服务
```mermaid
graph LR
A[开发者] --写好代码 --> B(本地编译成bin)
B -- dockerfile --> C[基于alpha镜像运行]
C -- docker-compose --> D[编排多个服务]
```
### GOLANG 服务目录树
<!-- apil_http
api2_http 
api3_http
grpc -->
<!-- 画个目录树 -->
```bash
├── api1_http
│   ├── main.go
├── api2_http
│   ├── main.go
├── api3_http
│   ├── main.go
├── grpc
│   ├── main.go
├── Makefile
|── docker-compose.yaml
``` 
四个golang服务,不用关心具体的代码,只需要知道他们是golang服务即可,3个api,1个grpc
### 编写makefile
```makefile
SERVICES = api1_http api2_http api3_http grpc
DOCKERFILE_CONTENT = FROM alpine:latest\nWORKDIR /app
# 定义 alpine:3.12 镜像为基础镜像
IMAGE = alpine:3.12
.PHONY: build dockerfiles deploy

build:
	@rm -rf ./bin
	@echo "Building $$@"
	@for service in $(SERVICES) ; do \
		CGO_ENABLED=0 GOOS=linux go build -o ./bin/$$service ./$$service ; \
	done

dockerfiles:
	@for service in $(SERVICES) ; do \
		mkdir -p ./dockerfiles/$$service ; \
		echo "FROM $(IMAGE)" > ./dockerfiles/$$service/Dockerfile ; \
		echo "WORKDIR /app" >> ./dockerfiles/$$service/Dockerfile ; \
		echo "COPY ./bin/$$service ./" >> ./dockerfiles/$$service/Dockerfile ; \
		echo "CMD [\"./$$service\"]" >> ./dockerfiles/$$service/Dockerfile ; \
	done

deploy: build dockerfiles
	@docker-compose up --build
```
makefile 文件内容包含了三个部分
1. build: 批量编译golang服务,生成二进制文件
2. dockerfiles: 批量生成dockerfile文件,基于alpine:3.12镜像,运行二进制文件
> 为什么要基于alpine镜像运行二进制文件,因为alpine镜像体积小,适合作为基础镜像
> 生成的dockerfile文件内容如下
···dockerfile
FROM alpine:3.12
WORKDIR /app
COPY ./bin/api1_http ./
CMD ["./api1_http"]
···
3. deploy: 通过docker-compose编排运行多个服务

### 编写docker-compose.yaml
```yaml
version: '3.8'
services:
  api1_http:
    build: 
      context: .
      dockerfile: ./dockerfiles/api1_http/Dockerfile
    ports:
      - 8081:8081
    env_file:
      - .env
  api2_http:
    build:
      context: .
      dockerfile: ./dockerfiles/api2_http/Dockerfile
    ports:
      - 8082:8081
    env_file:
      - .env

  api3_http:
    build:
      context: .
      dockerfile: ./dockerfiles/api3_http/Dockerfile
    ports:
      - 8083:8081
    env_file:
      - .env
  grpc:
    build:
      context: .
      dockerfile: ./dockerfiles/grpc/Dockerfile
    ports:
      - 8084:8082
    env_file:
      - .env
  jaeger:
    image: jaegertracing/all-in-one:1.6
    ports:
      - 16686:16686
      - 14268:14268
      - 16685:16685
```
docker-compose.yaml 包含了golang 服务 和 jaeger 服务
1. golang服务 的context 是当前目录, 所以dockerfile 文件里的COPY ./bin/api1_http ./ 会将当前目录下的bin/api1_http 复制到镜像的/app目录下
2. 这里还有一个.env文件,里面包含了环境变量,主要是在golang业务代码里使用,这里不做过多介绍
3. golang 业务里 会进行服务之间调用,集成了 jaeger 服务,所以需要将 jaeger 服务也编排进来
### 启动所有服务
```bash
make deploy
```
 会发现目录下多了一个bin目录,里面包含了编译好的二进制文件
 还有一个dockerfiles目录,里面包含了编译好的dockerfile文件 

### 查看 jaeger
我的 golang 服务 业务代码会互相调用,启动后 调用一个接口,就会在 jaeger 生成完整链路追踪
```
curl http://localhost:8081/api1
```
打开 jaeger 地址 http://localhost:16686
