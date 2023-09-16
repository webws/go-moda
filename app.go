package go_moda

import (
	"context"
	"errors"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/webws/go-moda/logger"
	"github.com/webws/go-moda/transport"
	modagrpc "github.com/webws/go-moda/transport/grpc"
	"github.com/webws/go-moda/transport/http"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"golang.org/x/sync/errgroup"

	"google.golang.org/grpc"
)

type (
	Option  func(o *options)
	options struct {
		name    string
		version string

		ctx        context.Context
		tracing    bool
		servers    []transport.Server
		httpServer http.HTTPServer
		grpcServer *modagrpc.Server
	}
)

func Name(name string) Option {
	return func(o *options) { o.name = name }
}

func Version(version string) Option {
	return func(o *options) { o.version = version }
}

func Tracing(Tracing bool) Option {
	return func(o *options) { o.tracing = Tracing }
}

func Server(srv ...transport.Server) Option {
	return func(o *options) { o.servers = srv }
}

type App struct {
	opts   options
	ctx    context.Context
	cancel func()
}

func New(opt ...Option) *App {
	options := options{
		ctx: context.Background(),
	}
	for _, o := range opt {
		o(&options)
	}
	ctx, cancel := context.WithCancel(options.ctx)
	app := &App{
		ctx:    ctx,
		cancel: cancel,
		opts:   options,
	}
	return app
}

func NewServer() *App {
	options := options{
		ctx: context.Background(),
	}
	ctx, cancel := context.WithCancel(options.ctx)
	app := &App{
		ctx:    ctx,
		cancel: cancel,
		opts:   options,
	}
	return app
}

func (a *App) SetTracing(tracing bool) *App {
	a.opts.tracing = tracing
	return a
}

func (a *App) AddHttpServer(address string, registerFunc func(g *gin.Engine)) *App {
	gin, httpServer := http.NewGinHttpServer(http.WithAddress(address))
	a.opts.servers = append(a.opts.servers, httpServer)
	registerFunc(gin)
	return a
}

func (a *App) AddGrpcServer(address string, registerFunc func(grpc.ServiceRegistrar)) *App {
	grpcServer := modagrpc.NewServer(modagrpc.WithServerAddress(address), modagrpc.WithRegisterFunc(registerFunc))
	a.opts.servers = append(a.opts.servers, grpcServer)
	return a
}

func (a *App) init() *App {
	if a == nil {
		a = NewServer()
	}
	if a.opts.httpServer != nil {
		if a.opts.tracing {
			a.opts.httpServer.EnableTracing()
		}
	}
	if a.opts.grpcServer != nil {
		a.opts.grpcServer.Server = grpc.NewServer()
		var grpcOption []grpc.ServerOption
		if a.opts.tracing {
			grpcOption = append(grpcOption, grpc.UnaryInterceptor(otelgrpc.UnaryServerInterceptor()))
			grpcOption = append(grpcOption, grpc.StreamInterceptor(otelgrpc.StreamServerInterceptor()))
			a.opts.grpcServer.Server = grpc.NewServer(grpcOption...)
		}
		a.opts.grpcServer.CallBackRegiser()

	}
	return a
}

func (a *App) Run() error {
	a = a.init()
	eg, ctx := errgroup.WithContext(a.ctx)
	for _, srv := range a.opts.servers {
		srv := srv
		eg.Go(func() error {
			<-ctx.Done() // wait app stop signal
			// logger.Infow("app run:server stop for stop signal")
			// srv.Stop(ctx) 后,start的那个协程解除阻塞,与这个ctx无关
			return srv.Stop(ctx)
		})
		eg.Go(func() error {
			return srv.Start(ctx)
		})
	}
	c := make(chan os.Signal, 1)
	signals := []os.Signal{syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT}
	signal.Notify(c, signals...)
	eg.Go(func() error {
		for {
			select {
			case <-ctx.Done():
				logger.Infow("app run: context done")
				return ctx.Err()
			case s := <-c:
				logger.Infow("app run: signal received", "signal", s.String())
				// 优雅退出,粗暴的可直接os.Exit(1)
				return a.Stop()
			}
		}
	})
	logger.Infow("app run:servers started")
	// eg.Wait() 等待协程优雅执行完毕
	if err := eg.Wait(); err != nil && !errors.Is(err, context.Canceled) {
		logger.Errorw("app run error", "error", err)
		return err
	}
	// app 结束
	logger.Infow("app run: servers stopped")
	return nil
}

func (a *App) Stop() error {
	if a.cancel != nil {
		a.cancel()
	}
	return nil
}
