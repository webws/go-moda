package config

import (
	"github.com/webws/go-moda/logger"
)

type Config struct {
	options options
}

/*
*
TODO 配置不符合toml
TODO 配置文件不存在
TODO new load unit test
*/
func New(opts ...Option) *Config {
	optionsObj := options{}
	for _, o := range opts {
		o(&optionsObj)
	}
	return &Config{options: optionsObj}
}

func (c *Config) Load(v interface{}) error {
	// 获取配置内容
	content, err := c.getSourceContent()
	if err != nil {
		logger.Errorw("Load.getSourceContent fail", "err", err)
		return err
	}
	if err := c.unmarshal(content, v); err != nil {
		logger.Errorw("Load.unmarshal fail", "err", err)
		return err
	}
	return nil
}

func (c *Config) getSourceContent() ([]byte, error) {
	var err error
	var content []byte
	if len(c.options.Sources) == 0 {
		// err := errors.New("error no config source usable")
		logger.Errorw("error no config source usable")
		return nil, nil
	}
	for _, source := range c.options.Sources {
		content, err = source.GetContent()
		if err == nil {
			break
		} else {
			logger.Errorw("getSourceContent.GetContent fail", "err", err, "sourceType", source.GetSourceName())
		}
	}
	return content, nil
}

// unmarshal 支持 toml.unmarshal
func (c *Config) unmarshal(p []byte, v interface{}) error {
	if c.options.unmarshaler == nil {
		c.options.unmarshaler = TomlUnmarshaler
	}
	c.options.unmarshaler(p, v)
	return nil
}
