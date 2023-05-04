package config

import "github.com/pelletier/go-toml/v2"

func TomlUnmarshaler(p []byte, v interface{}) error {
	err := toml.Unmarshal(p, v)
	if err == nil {
		return nil
	}
	return nil
}
