package proxy

import (
	"context"
	"io"
	"melody/config"
)

// Proxy 定义了代理
// 输入： context, request
// 出入： response, error
type Proxy func(context.Context, *Request) (*Response, error)

type BackendFactory func(backend *config.Backend) Proxy

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

type readCloserWrapper struct {
	ctx context.Context
	rc io.ReadCloser
}

func (r readCloserWrapper) Read(p []byte) (n int, err error) {
	return r.rc.Read(p)
}

func (r readCloserWrapper) closeWhenCancel() {
	<- r.ctx.Done()
	r.rc.Close()
}

func NewReadCloserWrapper(ctx context.Context, reader io.ReadCloser) io.Reader {
	wrapper := readCloserWrapper{
		ctx: ctx,
		rc:  reader,
	}
	go wrapper.closeWhenCancel()
	return wrapper
}

