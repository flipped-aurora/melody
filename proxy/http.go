package proxy

import (
	"context"
	"fmt"
	"melody/config"
	"melody/encoding"
	"melody/transport/http/client"
	"net/http"
	"strconv"
)

// HTTPResponseParser 将http.Response -> proxy.Response
type HTTPResponseParser func(context.Context, *http.Response) (*Response, error)

func HTTPProxyFactory(client *http.Client) BackendFactory {
	return CustomHTTPProxyFactory(func(_ context.Context) *http.Client { return client })
}

func CustomHTTPProxyFactory(cf client.HTTPClientFactory) BackendFactory {
	return func(backend *config.Backend) Proxy {
		return NewHTTPProxy(backend, cf, backend.Decoder)
	}
}

func NewHTTPProxy(remote *config.Backend, cf client.HTTPClientFactory, decode encoding.Decoder) Proxy {
	return NewHTTPProxyWithHTTPRequestExecutor(remote, client.DefaultHTTPRequestExecutor(cf), decode)
}

// NewHTTProxyWithHTTPRequestExecutor 将HTTPRequestExecutor封装，返回一个Proxy实例
func NewHTTPProxyWithHTTPRequestExecutor(remote *config.Backend, executor client.HTTPRequestExecutor, decoder encoding.Decoder) Proxy {
	if remote.Encoding == encoding.NOOP {
		return NewHTTPProxyDetailed(remote, executor, client.NoOpHTTPStatusHandler, NoOpHTTPResponseParser)
	}

	formatter := NewEntityFormatter(remote)
	responseParser := DefaultHTTPResponseParserFactory(HTTPResponseParserConfig{
		Decoder:         decoder,
		EntityFormatter: formatter,
	})
	return NewHTTPProxyDetailed(remote, executor, client.GetHTTPStatusHandler(remote), responseParser)
}

func NewHTTPProxyDetailed(remote *config.Backend, re client.HTTPRequestExecutor, ch client.HTTPStatusHandler, rp HTTPResponseParser) Proxy {
	return func(ctx context.Context, request *Request) (*Response, error) {
		requestToBackend, err := http.NewRequest(request.Method, request.URL.String(), request.Body)
		if err != nil {
			return nil, err
		}

		requestToBackend.Header = make(map[string][]string, len(request.Headers))

		for k, vs := range request.Headers {
			temp := make([]string, len(vs))
			copy(temp, vs)
			requestToBackend.Header[k] = temp
		}

		if request.Body != nil {
			if v, ok := request.Headers["Content-Length"]; ok && len(v) == 1 && v[0] != "chunked" {
				if size, err := strconv.Atoi(v[0]); err == nil {
					requestToBackend.ContentLength = int64(size)
				}
			}
		}

		// **真正发送请求的地方**
		resp ,err := re(ctx, requestToBackend)
		if requestToBackend.Body != nil {
			requestToBackend.Body.Close()
		}

		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		if err != nil {
			return nil, err
		}
		// response的成功或者错误处理
		resp, err = ch(ctx, resp)
		if err != nil {
			if t, ok := err.(responseError); ok {
				return &Response{
					Data: map[string]interface{}{
						fmt.Sprintf("error_%s", t.Name()): t,
					},
					Metadata: Metadata{StatusCode: t.StatusCode()},
				}, nil
			}
			return nil, err
		}

		return rp(ctx, resp)
	}
}

type responseError interface {
	Error() string
	Name() string
	StatusCode() int
}
