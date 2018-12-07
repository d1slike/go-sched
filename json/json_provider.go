package json

import "encoding/json"

var (
	Provider = &defaultProvider{}
)

type JsonProvider interface {
	Marshal(v interface{}) ([]byte, error)
	Unmarshal([]byte, interface{}) error
}

type defaultProvider struct {
}

func (p *defaultProvider) Marshal(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

func (p *defaultProvider) Unmarshal(bytes []byte, ptr interface{}) error {
	return json.Unmarshal(bytes, ptr)
}
