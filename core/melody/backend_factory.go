package melody

import (
	"context"
	"melody/config"
	"melody/logging"
	circuitbreaker "melody/middleware/melody-circuitbreaker/proxy"
	metrics "melody/middleware/melody-metrics/gin"
	"melody/proxy"
	"melody/transport/http/client"
)

// NewBackendFactory 创建BackendFactory，实际去请求每一个backend
func NewBackendFactoryWithContext(ctx context.Context, logger logging.Logger, metrics *metrics.Metrics) proxy.BackendFactory {
	clientFactory := client.NewHTTPClient
	httpRequestExecutor := client.DefaultHTTPRequestExecutor(clientFactory)
	backendFactory := func(backend *config.Backend) proxy.Proxy {
		return proxy.NewHTTPProxyWithHTTPRequestExecutor(backend, httpRequestExecutor, backend.Decoder)
	}
	// 使用断路器
	backendFactory = circuitbreaker.BackendFactory(backendFactory, logger)
	return backendFactory
}
