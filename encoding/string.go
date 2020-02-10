package encoding

import (
	"io"
	"io/ioutil"
)

const STRING = "string"

func NewStringDecoder(_ bool) func(io.Reader, *map[string]interface{}) error {
	return StringDecoder()
}

func StringDecoder() Decoder {
	return func(reader io.Reader, i *map[string]interface{}) error {
		data, err := ioutil.ReadAll(reader)
		if err != nil {
			return err
		}
		*i = map[string]interface{}{"content": string(data)}
		return nil
	}
}
