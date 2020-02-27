package bloomfilter

import "errors"

const Namespace = "melody_bloomfilter"

var (
	errNoConfig    = errors.New("no config for the bloomfilter")
	errWrongConfig = errors.New("invalid config for the bloomfilter")
)
