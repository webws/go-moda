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

