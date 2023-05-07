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

type HTTPServer interface {
	ServeHTTP(w http.ResponseWriter, r *http.Request)
	Start(address string) error
	Stop(ctx context.Context) error
	PprofRegister(string)
}

// 检查是否实现了HTTPServer接口
var _ HTTPServer = (*modagin.GinServer)(nil)

func NewGinServer() *modagin.GinServer {
	ginEngine := gin.Default()
	handle := &http.Server{
		Handler: ginEngine,
	}
	return &modagin.GinServer{
		Server: handle,
	}
}

// 检查是否实现了HTTPServer接口
var _ HTTPServer = (*modaecho.EchoServer)(nil)

func NewEchoServer() *modaecho.EchoServer {
	e := echo.New()
	e.HideBanner = true
	e.HidePort = true
	return &modaecho.EchoServer{Server: e}
}

// 检查是否实现了HTTPServer接口
var _ HTTPServer = (*netHttp.NetHTTPServer)(nil)

func NewNetHTTPServer() *netHttp.NetHTTPServer {
	serveHandle := http.NewServeMux()
	return &netHttp.NetHTTPServer{Server: &http.Server{
		Addr:    ":8081",
		Handler: serveHandle,
	}}
}
