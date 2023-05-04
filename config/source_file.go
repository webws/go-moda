package config

import (
	"io"
	"os"

	"github.com/webws/go-moda/logger"
)

type SourceFile struct {
	ConfigPath        string
	DefaultConfigPath string
}

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

func (sf *SourceFile) GetSourceName() string {
	return "File"
}
