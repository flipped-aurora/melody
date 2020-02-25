package melody

import (
	"context"
	"melody/config"
	"melody/logging"
	metrics "melody/middleware/melody-metrics/gin"
	"melody/proxy"
	"melody/transport/http/client"
)

// NewBackendFactory 创建BackendFactory，实际去请求每一个backend
func NewBackendFactoryWithContext(ctx context.Context, logger logging.Logger, metrics *metrics.Metrics) proxy.BackendFactory {
	clientFactory := client.NewHTTPClient
	httpRequestExecutor := client.DefaultHTTPRequestExecutor(clientFactory)
	return func(backend *config.Backend) proxy.Proxy {
		return proxy.NewHTTPProxyWithHTTPRequestExecutor(backend, httpRequestExecutor, backend.Decoder)
	}
}
