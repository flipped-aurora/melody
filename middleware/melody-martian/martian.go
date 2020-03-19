package martian

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"

	"melody/config"
	"melody/logging"
	"melody/proxy"
	"melody/transport/http/client"

	"github.com/google/martian"
	_ "github.com/google/martian/body"
	_ "github.com/google/martian/cookie"
	_ "github.com/google/martian/fifo"
	_ "github.com/google/martian/header"
	_ "github.com/google/martian/martianurl"
	"github.com/google/martian/parse"
	_ "github.com/google/martian/port"
	_ "github.com/google/martian/priority"
	_ "github.com/google/martian/stash"
	_ "github.com/google/martian/status"
)

// NewBackendFactory 创建一个代理。BackendFactory使用martian请求执行器包装注入的请求。如果解析额外的配置数据有任何问题，它只使用注入的请求执行程序。
func NewBackendFactory(logger logging.Logger, re client.HTTPRequestExecutor) proxy.BackendFactory {
	return NewConfiguredBackendFactory(logger, func(_ *config.Backend) client.HTTPRequestExecutor { return re })
}

// NewConfiguredBackendFactory ...
func NewConfiguredBackendFactory(logger logging.Logger, ref func(*config.Backend) client.HTTPRequestExecutor) proxy.BackendFactory {
	parse.Register("static.Modifier", staticModifierFromJSON)

	return func(remote *config.Backend) proxy.Proxy {
		re := ref(remote)
		result, ok := ConfigGetter(remote.ExtraConfig).(Result)
		if !ok {
			return proxy.NewHTTPProxyWithHTTPRequestExecutor(remote, re, remote.Decoder)
		}
		switch result.Err {
		case nil:
			return proxy.NewHTTPProxyWithHTTPRequestExecutor(remote, HTTPRequestExecutor(result.Result, re), remote.Decoder)
		case ErrEmptyValue:
			return proxy.NewHTTPProxyWithHTTPRequestExecutor(remote, re, remote.Decoder)
		default:
			logger.Error(result, remote.ExtraConfig)
			return proxy.NewHTTPProxyWithHTTPRequestExecutor(remote, re, remote.Decoder)
		}
	}
}

// HTTPRequestExecutor 在接收到的请求执行器上创建一个包装器，这样就可以在请求执行之前和之后执行martian modifiers
func HTTPRequestExecutor(result *parse.Result, re client.HTTPRequestExecutor) client.HTTPRequestExecutor {
	return func(ctx context.Context, req *http.Request) (resp *http.Response, err error) {
		if err = modifyRequest(result.RequestModifier(), req); err != nil {
			return
		}

		mctx, ok := req.Context().(*Context)
		if !ok || !mctx.SkippingRoundTrip() {
			resp, err = re(ctx, req)
			if err != nil {
				return
			}
			if resp == nil {
				err = ErrEmptyResponse
				return
			}
		} else if resp == nil {
			resp = &http.Response{
				Request:    req,
				Header:     http.Header{},
				StatusCode: http.StatusOK,
				Body:       ioutil.NopCloser(bytes.NewBufferString("")),
			}
		}

		err = modifyResponse(result.ResponseModifier(), resp)
		return
	}
}

func modifyRequest(mod martian.RequestModifier, req *http.Request) error {
	if req.Body == nil {
		req.Body = ioutil.NopCloser(bytes.NewBufferString(""))
	}
	if req.Header == nil {
		req.Header = http.Header{}
	}

	if mod == nil {
		return nil
	}
	return mod.ModifyRequest(req)
}

func modifyResponse(mod martian.ResponseModifier, resp *http.Response) error {
	if resp.Body == nil {
		resp.Body = ioutil.NopCloser(bytes.NewBufferString(""))
	}
	if resp.Header == nil {
		resp.Header = http.Header{}
	}
	if resp.StatusCode == 0 {
		resp.StatusCode = http.StatusOK
	}

	if mod == nil {
		return nil
	}
	return mod.ModifyResponse(resp)
}

// Namespace ...
const Namespace = "melody_martian"

// Result struct
type Result struct {
	Result *parse.Result
	Err    error
}

// ConfigGetter 实现 config.ConfigGetter interface.
func ConfigGetter(e config.ExtraConfig) interface{} {
	cfg, ok := e[Namespace]
	if !ok {
		return Result{nil, ErrEmptyValue}
	}

	data, ok := cfg.(map[string]interface{})
	if !ok {
		return Result{nil, ErrBadValue}
	}

	raw, err := json.Marshal(data)
	if err != nil {
		return Result{nil, ErrMarshallingValue}
	}

	r, err := parse.FromJSON(raw)

	return Result{r, err}
}

var (
	// ErrEmptyValue 是在名称空间中没有配置时返回的错误
	ErrEmptyValue = errors.New("getting the extra config for the martian module")
	// ErrBadValue 是配置不是映射时返回的错误
	ErrBadValue = errors.New("casting the extra config for the martian module")
	// ErrMarshallingValue 是配置映射不能再次编组时返回的错误
	ErrMarshallingValue = errors.New("marshalling the extra config for the martian module")
	// ErrEmptyResponse 是修改器接收到nil响应时返回的错误
	ErrEmptyResponse = errors.New("getting the http response from the request executor")
)
