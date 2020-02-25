package proxy

import (
	"context"
	"errors"
	"io"
	"melody/config"
)

// Namespace to be used in extra config
const Namespace = "melody_proxy"

var (
	// ErrNoBackends is the error returned when an endpoint has no backends defined
	ErrNoBackends = errors.New("all endpoints must have at least one backend")
	// ErrTooManyBackends is the error returned when an endpoint has too many backends defined
	ErrTooManyBackends = errors.New("too many backends for this proxy")
	// ErrTooManyProxies is the error returned when a middleware has too many proxies defined
	ErrTooManyProxies = errors.New("too many proxies for this proxy middleware")
	// ErrNotEnoughProxies is the error returned when an endpoint has not enough proxies defined
	ErrNotEnoughProxies = errors.New("not enough proxies for this endpoint")
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

type Middleware func(...Proxy) Proxy

func EmptyMiddleware(next ...Proxy) Proxy {
	if len(next) > 1 {
		panic(ErrTooManyProxies)
	}
	return next[0]
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

func NoopProxy(_ context.Context, _ *Request) (*Response, error) { return nil, nil }

