package main

import (
	"context"

	pbexample "github.com/webws/go-moda/example/pb/example"
)

// 实现 GrpcServer 接口
type ExampleServer struct {
	pbexample.UnimplementedExampleServiceServer
}

// 实现 GrpcServer 接口中的 SayHello 方法

func (s *ExampleServer) SayHello(ctx context.Context, req *pbexample.HelloRequest) (*pbexample.HelloResponse, error) {
	// return nil, nil
	return &pbexample.HelloResponse{Message: "Hello " + req.Name}, nil
}

// func (s *ExampleServer) SayHello(a context.Context, b *pbexample.HelloRequest) (*pbexample.HelloResponse, error) {

// 	return nil, nil
// 	// implementation
// }
