package xml

import (
	"github.com/clbanning/mxj"
	"io"
	"melody/encoding"
)

const XML = "xml"

type xmlReader struct {
	r io.Reader
}

func (x xmlReader) Read(p []byte) (n int, err error) {
	n, err = x.r.Read(p)

	if err != io.EOF {
		return n, err
	}

	if len(p) == n {
		return n, nil
	}

	p[n] = ([]byte("\n"))[0]
	return n + 1, err
}

func Register() {
	encoding.Register(XML, NewXMLDecoder)
}
func NewXMLDecoder(isCollection bool) func(io.Reader, *map[string]interface{}) error {
	if isCollection {
		return CollectionDecoder()
	} else {
		return Decoder()
	}
}

func CollectionDecoder() encoding.Decoder {
	return func(reader io.Reader, i *map[string]interface{}) error {
		mv, err := mxj.NewMapXmlReader(xmlReader{r: reader})
		if err != nil {
			return err
		}
		*i = map[string]interface{}{"collection": mv}
		return nil
	}
}

func Decoder() encoding.Decoder {
	return func(reader io.Reader, i *map[string]interface{}) error {
		mp, err := mxj.NewMapXmlReader(xmlReader{r: reader})
		if err != nil {
			return err
		}
		*i = mp
		return nil
	}
}
