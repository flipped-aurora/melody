package proxy

import (
	"io"
	"net/url"
)

// Request 包含了从Endpoint发送
// 到backends的数据
type Request struct {
	Method  string
	URL     *url.URL
	Query   url.Values
	Path    string
	Body    io.ReadCloser
	Params  map[string]string
	Headers map[string]string
}
