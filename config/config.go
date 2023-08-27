package config

import (
	"github.com/webws/go-moda/logger"
)

type (
	Option  func(*options)
	options struct {
		Sources []Source // 配置source
	}
)

func WithSources(sources []Source) Option {
	return func(o *options) {
		o.Sources = sources
	}
}

type Config struct {
	options options
}

func New(opts ...Option) *Config {
	optionsObj := options{}
	for _, o := range opts {
		o(&optionsObj)
	}
	return &Config{options: optionsObj}
}

func (c *Config) Load(v interface{}) error {
	for _, source := range c.options.Sources {
		if err := source.Unmarshal(v); err == nil {
			break
		} else {
			logger.Errorw("source.Unmarshal fail", "sourceType", source.GetSourceName(), "err", err)
		}
	}
	return nil
}
