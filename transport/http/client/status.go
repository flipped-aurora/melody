package client

import (
	"bytes"
	"context"
	"errors"
	"io/ioutil"
	"melody/config"
	"net/http"
)
const Namespace = "melody_http"

var ErrInvalidStatusCode = errors.New("Invalid status code")

// HTTPStatusHandler 将接受到的response中status code格式化
type HTTPStatusHandler func(context.Context, *http.Response) (*http.Response, error)

// HTTPResponseError 在某个Backend发生错误时，将封装进该对象的实例
type HTTPResponseError struct {
	Code int    `json:"http_status_code"`
	Msg  string `json:"http_body,omitempty"`
	name string
}


// NoOpHTTPStatusHandler 空实现
func NoOpHTTPStatusHandler(_ context.Context, resp *http.Response) (*http.Response, error) {
	return resp, nil
}

func GetHTTPStatusHandler(remote *config.Backend) HTTPStatusHandler {
	if e, ok := remote.ExtraConfig[Namespace]; ok {
		if m, ok := e.(map[string]interface{}); ok {
			if v, ok := m["return_error_details"]; ok {
				if b, ok := v.(string); ok && b != "" {
					return DetailedHTTPStatusHandler(DefaultHTTPStatusHandler, b)
				}
			}
		}
	}
	return DefaultHTTPStatusHandler
}

// DetailedHTTPStatusHandler
func DetailedHTTPStatusHandler(next HTTPStatusHandler, name string) HTTPStatusHandler {
	return func(ctx context.Context, resp *http.Response) (*http.Response, error) {
		if r, err := next(ctx, resp); err == nil {
			return r, nil
		}

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			body = []byte{}
		}
		resp.Body.Close()
		resp.Body = ioutil.NopCloser(bytes.NewBuffer(body))

		return resp, HTTPResponseError{
			Code: resp.StatusCode,
			Msg:  string(body),
			name: name,
		}
	}
}

func DefaultHTTPStatusHandler(ctx context.Context, resp *http.Response) (*http.Response, error) {
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, ErrInvalidStatusCode
	}

	return resp, nil
}


// Error returns the error message
func (r HTTPResponseError) Error() string {
	return r.Msg
}

// Name returns the name of the error
func (r HTTPResponseError) Name() string {
	return r.name
}

// StatusCode returns the status code returned by the backend
func (r HTTPResponseError) StatusCode() int {
	return r.Code
}

