package proxy

import (
	"context"
	"melody/config"
	"melody/logging"
	gobreaker "melody/middleware/melody-circuitbreaker"
	"melody/proxy"
)

func BackendFactory(next proxy.BackendFactory, logger logging.Logger) proxy.BackendFactory {
	return func(backend *config.Backend) proxy.Proxy {
		return NewMiddleware(backend, logger)(next(backend))
	}
}

func NewMiddleware(remote *config.Backend, logger logging.Logger) proxy.Middleware {
	config := gobreaker.ConfigGetter(remote.ExtraConfig).(gobreaker.Config)
	if config == gobreaker.DefaultCfg {
		return proxy.EmptyMiddleware
	}

	breaker := gobreaker.NewCircuitBreaker(config, logger)

	return func(next ...proxy.Proxy) proxy.Proxy {
		if len(next) > 1 {
			panic(proxy.ErrTooManyProxies)
		}

		return func(ctx context.Context, request *proxy.Request) (response *proxy.Response, err error) {
			res, err := breaker.Execute(func() (i interface{}, err error) {
				return next[0](ctx, request)
			})
			if err != nil {
				return nil, err
			}
			return res.(*proxy.Response), err
		}
	}
}
