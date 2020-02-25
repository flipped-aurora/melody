package melody

import (
	"melody/logging"
	"melody/proxy"
)

func NewProxyFactory(logger logging.Logger, backend proxy.BackendFactory) proxy.Factory {
	// 完成了默认的ProxyFactory
	// TODO 与其他服务集成
	return proxy.NewDefaultFactory(backend, logger)
}
