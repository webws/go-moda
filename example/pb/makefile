# 转换exlmple 为grpc go文件
# TODO shell batch 
Example2Go:
	protoc --go-grpc_out=. --go_opt=paths=source_relative  --go-grpc_opt=paths=source_relative --go_out=. example/*.proto
#install proto proto-gen-go
Pb2Go:
	protoc --go-grpc_out=. --go_opt=paths=source_relative  --go-grpc_opt=paths=source_relative --go_out=. *.proto
#install proto proto-gen-go
Install:
	go get -u google.golang.org/protobuf/cmd/protoc-gen-go
	go get -u google.golang.org/grpc/cmd/protoc-gen-go-grpc

Gen:
	make Install && make Pb2Go