package grpc

import (
	"context"

	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type dialOption func(*dialOptions)

// grpc clinet options
type dialOptions struct {
	grpcOptions  []grpc.DialOption
	tracing      bool
	withInsecure bool
}

// tracing
func WithDialTracing(tracing bool) dialOption {
	return func(o *dialOptions) {
		o.tracing = tracing
	}
}

// withInsecure
func WithDialWithInsecure(withInsecure bool) dialOption {
	return func(o *dialOptions) {
		o.withInsecure = withInsecure
	}
}

// grpc dial options
func WithDialOptions(opts ...grpc.DialOption) dialOption {
	return func(o *dialOptions) {
		o.grpcOptions = opts
	}
}

// Dial 封装了 grpc.DialContext
func Dial(ctx context.Context, target string, option ...dialOption) (*grpc.ClientConn, error) {
	opts := &dialOptions{}
	for _, o := range option {
		o(opts)
	}
	if opts.withInsecure {
		opts.grpcOptions = append(opts.grpcOptions, grpc.WithTransportCredentials(insecure.NewCredentials()))
		// opts.grpcOptions = append(opts.grpcOptions, grpc.WithInsecure())
	}
	if opts.tracing {
		// grpc 客户端 启用链路追踪
		opts.grpcOptions = append(opts.grpcOptions,
			grpc.WithUnaryInterceptor(otelgrpc.UnaryClientInterceptor()),
			grpc.WithStreamInterceptor(otelgrpc.StreamClientInterceptor()),
		)
	}
	return grpc.DialContext(ctx, target, opts.grpcOptions...)
}
