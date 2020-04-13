package encoding

import (
	"encoding/json"
	"io"
)

const JSON = "json"

func NewJSONDecoder(isCollection bool) func(io.Reader, *map[string]interface{}) error {
	if isCollection {
		return JSONCollectionDecoder()
	} else {
		return JSONDecoder()
	}
}

func JSONDecoder() Decoder {
	return func(reader io.Reader, i *map[string]interface{}) error {
		d := json.NewDecoder(reader)
		d.UseNumber()
		return d.Decode(i)
	}
}

func JSONCollectionDecoder() Decoder {
	return func(reader io.Reader, i *map[string]interface{}) error {
		var collection []interface{}
		d := json.NewDecoder(reader)
		d.UseNumber()
		if err := d.Decode(&collection); err != nil {
			return err
		}
		*i = map[string]interface{}{"collection": collection}
		return nil
	}
}
