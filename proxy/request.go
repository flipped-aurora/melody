package proxy

import (
	"bytes"
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
	Headers map[string][]string
}

func (r *Request) Clone() Request {
	return Request{
		Method:  r.Method,
		URL:     r.URL,
		Query:   r.Query,
		Path:    r.Path,
		Body:    r.Body,
		Params:  r.Params,
		Headers: r.Headers,
	}
}

// CloneRequest 返回一个请求副本
func CloneRequest(r *Request) *Request {
	clone := r.Clone()
	clone.Headers = CloneRequestHeaders(r.Headers)
	clone.Params = CloneRequestParams(r.Params)
	return &clone
}

// CloneRequestHeaders 返回一个接收到的请求头的副本
func CloneRequestHeaders(headers map[string][]string) map[string][]string {
	m := make(map[string][]string, len(headers))
	for k, vs := range headers {
		tmp := make([]string, len(vs))
		copy(tmp, vs)
		m[k] = tmp
	}
	return m
}

// CloneRequestParams 返回接收到的请求参数的副本
func CloneRequestParams(params map[string]string) map[string]string {
	m := make(map[string]string, len(params))
	for k, v := range params {
		m[k] = v
	}
	return m
}

func (r *Request) GeneratePath(URLPattern string) {
	if len(r.Params) == 0 {
		r.Path = URLPattern
		return
	}
	buff := []byte(URLPattern)
	for k, v := range r.Params {
		key := []byte{}
		key = append(key, "{{."...)
		key = append(key, k...)
		key = append(key, "}}"...)
		buff = bytes.Replace(buff, key, []byte(v), -1)
	}
	r.Path = string(buff)
}
