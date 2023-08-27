package config

import (
	"testing"

	"github.com/webws/go-moda/logger"
)

type ConfigTest struct {
	HttpAddr  string `json:"http_addr" toml:"http_addr" yaml:"http_addr"`
	GrpcAddr  string `json:"grpc_addr" toml:"grpc_addr" yaml:"grpc_addr"`
	JaegerUrl string `json:"jaeger_url" toml:"jaeger_url" yaml:"jaeger_url" mapstructure:"jaeger_url"`
	Tracing   bool   `toml:"tracing"  json:"tracing" yaml:"tracing" ` // opentelemetry tracing

	Hostname    string      `json:"hostname" toml:"hostname" yaml:"hostname"`
	Hostname2   string      `json:"host_name2" toml:"host_name2" yaml:"host_name2"`
	ServiceAddr ServiceAddr `json:"service_addr" toml:"service_addr"` // 依赖的服务地址
}
type ServiceAddr struct {
	Api1  string `json:"api1" toml:"api1" yaml:"api1"`
	Api2  string `json:"api2" toml:"api2" yaml:"api2"`
	Api3  string `json:"api3" toml:"api3" yaml:"api3"`
	Grpc1 string `json:"grpc1" toml:"grpc1" yaml:"grpc1"`
}

func TestSourceFile_Unmarshal2(t *testing.T) {
	//  定义测试结构提,一个字段 file_path 用于测试配置文件路径
	type test struct {
		FilePath string `json:"file_path" toml:"file_path" yaml:"file_path"`
	}
	// 测试数据 []test
	tests := []test{
		{FilePath: "./test_yaml.yaml"},
		{FilePath: "./test_toml.toml"},
	}
	for _, v := range tests {
		sf := &SourceFile{
			ConfigPath:        v.FilePath,
			DefaultConfigPath: v.FilePath,
		}
		c := &ConfigTest{}
		err := sf.Unmarshal(c)
		if err != nil {
			t.Error(err)
		}
		logger.Infow("Unmarshal file sucess", "source", sf.GetSourceName(), "v", c, "file_path", v.FilePath)
	}
}
