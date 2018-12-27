package internal

import (
	"github.com/d1slike/go-sched/json"
)

func castData(data interface{}) ([]byte, error) {
	switch d := data.(type) {
	case []byte:
		return d, nil
	case string:
		return []byte(d), nil
	default:
		return json.Provider.Marshal(data)
	}
}
