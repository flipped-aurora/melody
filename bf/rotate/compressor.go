package rotate

import (
	"compress/gzip"
	"io"
)

type Compressor interface {
	NewWriter(io.Writer) io.WriteCloser
	NewReader(io.Reader) (io.Reader, error)
}

var compressor Compressor = new(Gzip)

func SetCompressor(c Compressor) {
	compressor = c
}

type Gzip int

func (g *Gzip) NewWriter(w io.Writer) io.WriteCloser {
	return gzip.NewWriter(w)
}

func (g *Gzip) NewReader(r io.Reader) (io.Reader, error) {
	return gzip.NewReader(r)
}
