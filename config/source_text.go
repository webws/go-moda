package config

type SourceText struct {
	Context string
}

func (st *SourceText) GetContent() ([]byte, error) {
	return nil, nil
}
func (sf *SourceText) GetSourceName() string {
	return "test"
}
