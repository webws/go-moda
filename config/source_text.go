package config

import (
	"github.com/pelletier/go-toml/v2"
	"github.com/webws/go-moda/logger"
)

var _ Source = (*SourceText)(nil)

type SourceText struct {
	Context string
}

// TODO: support more format, json, yaml
func (st *SourceText) Unmarshal(v interface{}) error {
	if err := toml.Unmarshal([]byte(st.Context), v); err != nil {
		logger.Errorw("Unmarshal fail", "err", err)
		return err
	}
	return nil
}

func (st *SourceText) GetSourceName() SourceName {
	return SourceNameText
}
