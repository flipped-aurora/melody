package encoding

import (
	"io"
	"melody/register"
)

var (
	decoders        = initDecoderRegister()
	defaultDecoders = map[string]func(bool) func(io.Reader, *map[string]interface{}) error{
		JSON:   NewJSONDecoder,
		STRING: NewStringDecoder,
		NOOP:   NoOpDecoderFactory,
	}
)

type DecoderRegister struct {
	data register.Untyped
}

func (d *DecoderRegister) Get(s string) DecoderFactory {
	// if can not get decoder:key = s,
	// return json decoder
	for _, v := range []string{s, JSON} {
		if v, ok := d.data.Get(v); ok {
			decoderFactory, ok := v.(func(bool) func(io.Reader, *map[string]interface{}) error)
			if ok{
				return decoderFactory
			}
		}
	}

	return NewJSONDecoder
}

func (d *DecoderRegister) Register(name string, factory func(bool) func(io.Reader, *map[string]interface{}) error) error {
	d.data.Register(name, factory)
	return nil
}

func initDecoderRegister() *DecoderRegister {
	decoder := &DecoderRegister{data: register.New()}
	for k, v := range defaultDecoders {
		decoder.data.Register(k, v)
	}
	return decoder
}
