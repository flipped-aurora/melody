package proxy

import (
	"context"
	"io"
)

// Proxy 定义了代理
// 输入： context, request
// 出入： response, error
type Proxy func(context.Context, *Request) (*Response, error)

// Metadata 包含了response header 和 response code
type Metadata struct {
	Headers    map[string][]string
	StatusCode int
}

// Response 作为proxy的输出
type Response struct {
	Data       map[string]interface{}
	IsComplete bool
	Io         io.Reader
	Metadata   Metadata
}
