package go_moda

import (
	"context"
	"sync"
	"testing"

	"github.com/webws/go-moda/logger"
)

var onceStart sync.Once

// 定义测试serer,继承transport.Server
type testServer struct {
	Ctx    context.Context
	Cancel context.CancelFunc
	Name   string
	Once   sync.Once
}

//  NewTestServer 初始化 ctx 和 cancel
func NewTestServer(name string) *testServer {
	ctx, cancel := context.WithCancel(context.Background())
	return &testServer{
		Ctx:    ctx,
		Cancel: cancel,
		Name:   name,
	}
}

// start 服务，这里是模拟http服务启动，需要一个程序一直运行监听http请求,另一个程序检测是否需要停止服务
func (s *testServer) Start(ctx context.Context) error {
	for {
		select {
		case <-s.Ctx.Done():
			// record log:start 程序里 检测到服务已关闭,返回
			logger.Infow("testServer is closed", "name", s.Name)
			return nil
		default:
			s.Once.Do(func() {
				logger.Infow("testServer start", "name", s.Name)
			})
		}
	}
}

// stop 服务，这里是模拟服务停止，调用cancel()方法
func (s *testServer) Stop(ctx context.Context) error {
	logger.Infow("testServer stop", "name", s.Name)
	s.Cancel()
	return nil
}

//  go test -v -run TestApp_Run
func TestApp_Run(t *testing.T) {
	a := New(
		Server(NewTestServer("server1"), NewTestServer("server2")),
		Version("1.0.0"),
	)
	a.Run()
}
