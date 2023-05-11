package http

import (
	"context"
	"net/http"
	"net/http/pprof"

	"github.com/labstack/echo/v4"
	"github.com/webws/go-moda/logger"
	"go.opentelemetry.io/contrib/instrumentation/github.com/labstack/echo/otelecho"
)

type EchoServer struct {
	Server *echo.Echo
}

func (e *EchoServer) GetServer() *echo.Echo {
	return e.Server
}

func (e *EchoServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	e.Server.ServeHTTP(w, r)
}

func (e *EchoServer) Start(address string) error {
	logger.Infow("EchoServer start", "address", address)
	return e.Server.Start(address)
}

func (e *EchoServer) Stop(ctx context.Context) error {
	return e.Server.Shutdown(ctx)
}

func (e *EchoServer) PprofRegister(pprofPrefix string) {
	g := e.GetServer().Group(pprofPrefix)
	g.GET("/*any", echo.WrapHandler(http.HandlerFunc(pprof.Index)))
}

func (e *EchoServer) EnableTracing() {
	logger.Infow("EchoServer EnableTracing")
	e.Server.Use(otelecho.Middleware("moda"))
}
