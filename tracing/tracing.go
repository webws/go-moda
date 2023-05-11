package tracing

import (
	"context"

	"go.opentelemetry.io/contrib/propagators/b3"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.7.0"
	"go.opentelemetry.io/otel/trace"

	// go-moda logger
	"github.com/webws/go-moda/logger"
)

var tracer = otel.Tracer("go-moda")

// 启动 jaeger ~/tool/jaeger-1.31.0-linux-amd64/jaeger-all-in-one
// 查看面板 地址 http://localhost:16686/
// 集成tracer Batcher 地址 http://localhost:14268/api/traces

func InitJaegerProvider(jaegerUrl string, serviceName string) (func(ctx context.Context) error, error) {
	if jaegerUrl == "" {
		logger.Errorw("jaeger url is empty")
		return nil, nil
	}
	tracer = otel.Tracer(serviceName)
	exp, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(jaegerUrl)))
	if err != nil {
		return nil, err
	}
	tp := tracesdk.NewTracerProvider(
		tracesdk.WithBatcher(exp),
		tracesdk.WithResource(resource.NewSchemaless(
			semconv.ServiceNameKey.String(serviceName),
		)),
	)
	otel.SetTracerProvider(tp)
	// otel.SetTextMapPropagator(propagation.TraceContext{})
	b3Propagator := b3.New(b3.WithInjectEncoding(b3.B3MultipleHeader))
	propagator := propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}, b3Propagator)
	otel.SetTextMapPropagator(propagator)
	return tp.Shutdown, nil
}

func Start(ctx context.Context, name string) (context.Context, trace.Span) {
	return tracer.Start(ctx, name)
}
