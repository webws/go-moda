package config

type Source interface {
	GetContent() ([]byte, error)
	GetSourceName() string
}
