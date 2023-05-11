package main

import (
	"context"
	"time"

	"github.com/webws/go-moda/tracing"
	"go.opentelemetry.io/otel/attribute"
)

func Bar(ctx context.Context) {
	// Use the global TracerProvider.
	ctx, span := tracing.Start(ctx, "service.bar")
	defer span.End()

	span.SetAttributes(attribute.Key("bar.attrbute").String("value"))
	time.Sleep(time.Second * 2)
	Bar2(ctx)
	defer span.End()
	// Do bar...
}

func Bar2(ctx context.Context) {
	// Use the global TracerProvider.
	ctx, span := tracing.Start(ctx, "service.bar2")
	defer span.End()
	span.SetAttributes(attribute.Key("bar2.attrbute").String("value"))
	time.Sleep(time.Second * 4)
	Bar3(ctx)
	defer span.End()

	// Do bar...
}

func Bar3(ctx context.Context) {
	// Use the global TracerProvider.
	_, span := tracing.Start(ctx, "service.bar3")
	defer span.End()
	span.SetAttributes(attribute.Key("bar3.attrbute").String("value"))
	time.Sleep(time.Second * 1)
	defer span.End()
}
