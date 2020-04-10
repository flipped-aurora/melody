package convert

import (
	"encoding/json"
	"time"
)

func ObjectToStringTime(v interface{}, format string) (string, bool) {
	if n, ok := v.(json.Number); ok {
		if ti, err := n.Int64(); err == nil {
			t := time.Unix(ti, 0).Format(format)
			return t, true
		}
	}
	return "", false
}

func ObjectToInt(v interface{}) (int64, bool) {
	if n, ok := v.(json.Number); ok {
		if ti, err := n.Int64(); err == nil {
			return ti, true
		}
	}
	return 0, false
}

func ObjectToFloat(v interface{}) (float64, bool) {
	if n, ok := v.(json.Number); ok {
		if ti, err := n.Float64(); err == nil {
			return ti, true
		}
	}
	return 0, false
}
