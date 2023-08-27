package config

// Source 配置来源
type SourceName string

const (
	SourceNameFile   SourceName = "File"
	SourceNameEtcd   SourceName = "Etcd"
	SourceNameConsul SourceName = "Consul"
	SourceNameText   SourceName = "Text"
)

type Source interface {
	GetSourceName() SourceName
	Unmarshal(v interface{}) error
}
