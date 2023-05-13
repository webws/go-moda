package config

import (
	"os"

	"github.com/webws/go-moda/logger"
)

var Conf = &Config{}

type Config struct {
	HttpAddr    string      `json:"http_addr" toml:"http_addr"`
	GrpcAddr    string      `json:"grpc_addr" toml:"grpc_addr"`
	JaegerUrl   string      `json:"jaeger_url" toml:"jaeger_url"`
	Tracing     bool        `toml:"tracing"  json:"tracing"`          // opentelemetry tracing
	ServiceAddr ServiceAddr `json:"service_addr" toml:"service_addr"` // 依赖的服务地址
}
type ServiceAddr struct {
	Api1  string `json:"api1" toml:"api1"`
	Api2  string `json:"api2" toml:"api2"`
	Api3  string `json:"api3" toml:"api3"`
	Grpc1 string `json:"grpc1" toml:"grpc1"`
}

// 获取环境变量配置, 优先级高于配置文件
func (c *Config) SetEnvServiceAddr() {
	// Getenv svc addr
	serviceApp1 := os.Getenv("SERVICE_APP1")
	serviceApp2 := os.Getenv("SERVICE_APP2")
	serviceApp3 := os.Getenv("SERVICE_APP3")
	serviceGrpc1 := os.Getenv("SERVICE_GRPC1")
	JaegerUrl := os.Getenv("JAEGER_URL")
	Tracing := os.Getenv("TRACING")

	if serviceApp1 != "" {
		c.ServiceAddr.Api1 = serviceApp1
	}
	if serviceApp2 != "" {
		c.ServiceAddr.Api2 = serviceApp2
	}
	if serviceApp3 != "" {
		c.ServiceAddr.Api3 = serviceApp3
	}
	if serviceGrpc1 != "" {
		c.ServiceAddr.Grpc1 = serviceGrpc1
	}
	if JaegerUrl != "" {
		c.JaegerUrl = JaegerUrl
	}
	if Tracing != "" {
		c.Tracing = true
	}
	// 打印 serveice addr
	logger.Infow("env service addr", "service_addr", c.ServiceAddr, "jaeger_url", c.JaegerUrl, "tracing", c.Tracing)
}
