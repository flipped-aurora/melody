package client

import (
	"context"
	"net/http"
)

var defaultHTTPClient = &http.Client{}

// HTTPRequestExecutor 定义了代理的client，由该方法去实现代理
type HTTPRequestExecutor func(context.Context, *http.Request) (*http.Response, error)

// HTTPClientFactory 根据context定制client
type HTTPClientFactory func(context.Context) *http.Client

func NewHTTPClient(context.Context) *http.Client {
	return defaultHTTPClient
}

// DeafultHTTPRequestExecutor 默认的request执行器
func DefaultHTTPRequestExecutor(clientFactory HTTPClientFactory) HTTPRequestExecutor {
	return func(i context.Context, request *http.Request) (response *http.Response, e error) {
		return clientFactory(i).Do(request)
	}
}
