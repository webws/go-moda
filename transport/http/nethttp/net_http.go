package http

import (
	"context"
	"fmt"
	"net/http"
	"net/http/pprof"

	"github.com/webws/go-moda/logger"
)

type NetHTTPServer struct {
	Server *http.Server
}

func NewNetHTTPServer() *NetHTTPServer {
	serveHandle := http.NewServeMux()
	return &NetHTTPServer{Server: &http.Server{
		Addr:    ":8081",
		Handler: serveHandle,
	}}
}

func (d *NetHTTPServer) GetServer() *http.ServeMux {
	return d.Server.Handler.(*http.ServeMux)
}

func (d *NetHTTPServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	d.Server.Handler.ServeHTTP(w, r)
}

func (d *NetHTTPServer) Start(address string) error {
	if address != "" {
		d.Server.Addr = address
	}
	logger.Infow("NetHTTPServer start", "address", address)
	return d.Server.ListenAndServe()
}

func (d *NetHTTPServer) Stop(ctx context.Context) error {
	return d.Server.Shutdown(ctx)
}

// pprof register
func (d *NetHTTPServer) PprofRegister(pprofPrefix string) {
	pattern := fmt.Sprintf("%s/", pprofPrefix)
	fmt.Println(pattern)
	d.GetServer().HandleFunc(pattern, http.HandlerFunc(pprof.Index))
}

func (d *NetHTTPServer) EnableTracing() {
	// TODO: gin middleware
}
