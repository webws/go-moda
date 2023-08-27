package config

import (
	"io"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
	"github.com/webws/go-moda/logger"

	"github.com/mitchellh/mapstructure"
)

var _ Source = (*SourceFile)(nil)

type SourceFile struct {
	ConfigPath        string
	DefaultConfigPath string
}

// GetSourceName implements Source
func (*SourceFile) GetSourceName() SourceName {
	return SourceNameFile
}

// GetContent Deprecated: use Unmarshal instead
func (sf *SourceFile) GetContent() ([]byte, error) {
	configPath := sf.ConfigPath
	if configPath == "" {
		configPath = sf.DefaultConfigPath
	}
	f, err := os.Open(configPath)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	content, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}
	logger.Infow("GetContent file sucess", "source", sf.GetSourceName())
	return content, nil
}

// Unmarshal implements Source unmarshal to struct
func (sf *SourceFile) Unmarshal(v interface{}) error {
	configPath := sf.DefaultConfigPath
	if sf.ConfigPath != "" {
		configPath = sf.ConfigPath
	}
	return ViperUnmarshal(v, configPath)
}

// ViperUnmarshal use viper unmarshal to struct
// auto set decode tag_name by configPath
func ViperUnmarshal(v interface{}, configPath string) error {
	var tagName string
	ext := filepath.Ext(configPath)
	if len(ext) > 1 {
		tagName = ext[1:]
	}
	// set decode tag_name, default is mapstructure
	decoderConfigOption := func(c *mapstructure.DecoderConfig) {
		c.TagName = tagName
	}
	cViper := viper.New()
	cViper.SetConfigFile(configPath)
	if err := cViper.ReadInConfig(); err != nil {
		return err
	}
	return cViper.Unmarshal(v, decoderConfigOption)
}
