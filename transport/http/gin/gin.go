package http

import (
	"context"
	"net/http"
	"net/http/pprof"

	"github.com/gin-gonic/gin"
	"github.com/webws/go-moda/logger"

	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

type GinServer struct {
	Server *http.Server
	// server *gin.Engine
}

func (g *GinServer) GetServer() *gin.Engine {
	return g.Server.Handler.(*gin.Engine)
}

func (g *GinServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	g.Server.Handler.ServeHTTP(w, r)
}

func (g *GinServer) Start(address string) error {
	g.Server.Addr = address
	logger.Infow("GinServer start", "address", address)
	if err := g.Server.ListenAndServe(); err != nil {
		return err
	}
	return nil
}

func (g *GinServer) Stop(ctx context.Context) error {
	return g.Server.Shutdown(ctx)
}

func (g *GinServer) PprofRegister(prefix string) {
	r := g.GetServer().Group(prefix)
	r.GET("/*any", gin.WrapF(pprof.Index))
}

func (g *GinServer) EnableTracing() {
	g.GetServer().Use(otelgin.Middleware("my-server"))
}
