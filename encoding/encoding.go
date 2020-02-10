package encoding

import "io"

const NOOP = "no-op"

// NoOpDecoder implements the Decoder interface
func NoOpDecoder(_ io.Reader, _ *map[string]interface{}) error { return nil }

func NoOpDecoderFactory(_ bool) func(io.Reader, *map[string]interface{}) error { return NoOpDecoder }

// Decoder a accept a param that can be read , and a target map pointer to write
type Decoder func(io.Reader, *map[string]interface{}) error

// DecoderFactory returns a Decoder
// 1. EntityDecoder {}
// 2. CollectionDecoder []
type DecoderFactory func(bool) func(io.Reader, *map[string]interface{}) error

func Get(key string) DecoderFactory {
	return decoders.Get(key)
}

func Register(key string, factory DecoderFactory) error {
	return decoders.Register(key, factory)
}
