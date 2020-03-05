package proxy

import (
	"melody/config"
	"melody/logging"
	"melody/sd"
)

// Factory 作为ProxyFactory的标准
type Factory interface {
	New(cfg *config.EndpointConfig) (Proxy, error)
}

type FactoryFunc func(*config.EndpointConfig) (Proxy, error)

// New implements the Factory interface
func (f FactoryFunc) New(cfg *config.EndpointConfig) (Proxy, error) { return f(cfg) }

type defaultFactory struct {
	backendFactory    BackendFactory
	logger            logging.Logger
	subscriberFactory sd.SubscriberFactory
}

func NewDefaultFactory(factory BackendFactory, logger logging.Logger) Factory {
	return NewDefaultFactoryWithSubscriberFactory(factory, logger, sd.GetSubscriber)
}

func NewDefaultFactoryWithSubscriberFactory(factory BackendFactory, logger logging.Logger, subscriberFactory sd.SubscriberFactory) Factory {
	return defaultFactory{
		backendFactory:    factory,
		logger:            logger,
		subscriberFactory: subscriberFactory,
	}
}

func (d defaultFactory) New(cfg *config.EndpointConfig) (p Proxy, err error) {
	switch len(cfg.Backends) {
	case 0:
		err = ErrNoBackends
	case 1:
		p, err = d.NewSingle(cfg)
	default:
		p, err = d.NewMulti(cfg)
	}

	if err != nil {
		return
	}
	// 执行顺序：⑥
	p = NewStaticDataMiddleware(cfg)(p)
	return
}

func (d defaultFactory) NewStack(backend *config.Backend) (p Proxy) {
	// 根据config.Backend定制backendProxy 执行顺序：④
	p = d.backendFactory(backend)
	// 均衡中间件注册                     执行顺序：③
	p = NewLoadBalancedMiddlewareWithSubscriber(d.subscriberFactory(backend))(p)
	if backend.ConcurrentCalls > 1 {
		// 并发调用 > 1                    执行顺序：②
		p = NewConcurrentCallMiddleware(backend)(p)
	}
	// 基础的Request构造器                 执行顺序：①
	p = NewRequestBuilderMiddleware(backend)(p)
	return
}

func (d defaultFactory) NewSingle(endpointConfig *config.EndpointConfig) (Proxy, error) {
	return d.NewStack(endpointConfig.Backends[0]), nil
}

func (d defaultFactory) NewMulti(endpointConfig *config.EndpointConfig) (p Proxy, err error) {

	backendProxies := make([]Proxy, len(endpointConfig.Backends))
	for i, v := range endpointConfig.Backends {
		backendProxies[i] = d.NewStack(v)
	}
	// 执行顺序：⑤
	p = NewMergeDataMiddleware(endpointConfig)(backendProxies...)
	return
}
