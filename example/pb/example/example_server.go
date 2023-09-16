package example

import (
	"context"
)

type ExampleServer struct {
	UnimplementedExampleServiceServer
}

func (s *ExampleServer) SayHello(ctx context.Context, req *HelloRequest) (*HelloResponse, error) {
	return &HelloResponse{Message: "Hello " + req.Name}, nil
}
