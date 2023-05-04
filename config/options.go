package config

type Unmarshaler func(p []byte, v interface{}) error

type Option func(*options)
type options struct {
	version     string
	name        string
	Sources     []Source    // 配置source
	unmarshaler Unmarshaler // 将配置内容解析到 app 定义的conf结构方法
}

func WithName(name string) Option {
	return func(o *options) {
		o.name = name
	}
}
func WithSources(sources []Source) Option {
	return func(o *options) {
		o.Sources = sources
	}
}
