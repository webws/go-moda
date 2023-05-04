package go_moda

import (
	"context"
	"errors"
	"os"
	"os/signal"
	"syscall"

	"github.com/webws/go-moda/logger"
	"github.com/webws/go-moda/transport"
	"golang.org/x/sync/errgroup"
)

type (
	Option  func(o *options)
	options struct {
		name    string
		version string

		ctx context.Context

		servers []transport.Server
	}
)

func Name(name string) Option {
	return func(o *options) { o.name = name }
}

func Version(version string) Option {
	return func(o *options) { o.version = version }
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

func (a *App) Run() error {
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
