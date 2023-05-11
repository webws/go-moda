package http

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/labstack/echo/v4"
	modaecho "github.com/webws/go-moda/transport/http/echo"
	modagin "github.com/webws/go-moda/transport/http/gin"
	netHttp "github.com/webws/go-moda/transport/http/nethttp"
)

var (
	_ HTTPServer = (*modagin.GinServer)(nil)
	_ HTTPServer = (*modaecho.EchoServer)(nil)
	_ HTTPServer = (*netHttp.NetHTTPServer)(nil)
)

type HTTPServer interface {
	ServeHTTP(w http.ResponseWriter, r *http.Request)
	Start(address string) error
	Stop(ctx context.Context) error
	PprofRegister(string)
	EnableTracing()
}

func newGinServer() *modagin.GinServer {
	ginEngine := gin.Default()
	handle := &http.Server{
		Handler: ginEngine,
	}
	return &modagin.GinServer{
		Server: handle,
	}
}

func newEchoServer() *modaecho.EchoServer {
	e := echo.New()
	e.HideBanner = true
	e.HidePort = true
	return &modaecho.EchoServer{Server: e}
}

func newNetHTTPServer() *netHttp.NetHTTPServer {
	serveHandle := http.NewServeMux()
	return &netHttp.NetHTTPServer{Server: &http.Server{
		Handler: serveHandle,
	}}
}
