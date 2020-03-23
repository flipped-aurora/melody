package melody

import (
	"melody/logging"
	jsonschema "melody/middleware/melody-jsonschema"
	metrics "melody/middleware/melody-metrics/gin"
	"melody/proxy"
)

func NewProxyFactory(logger logging.Logger, backend proxy.BackendFactory, metrics *metrics.Metrics) proxy.Factory {
	// 完成了默认的ProxyFactory
	proxyFactory := proxy.NewDefaultFactory(backend, logger)
	proxyFactory = proxy.NewShadowFactory(proxyFactory)
	proxyFactory = jsonschema.ProxyFactory(proxyFactory)
	proxyFactory = metrics.NewProxyFactory("endpoint", proxyFactory)
	return proxyFactory

}
