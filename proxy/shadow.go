package proxy

import (
	"context"
	"melody/config"
)

const (
	shadowKey = "shadow"
)

type shadowFactory struct {
	f Factory
}

// New check the Backends for an ExtraConfig with the "shadow" param to true
// 实现工厂接口
func (s shadowFactory) New(cfg *config.EndpointConfig) (p Proxy, err error) {
	if len(cfg.Backends) == 0 {
		err = ErrNoBackends
		return
	}

	shadow := []*config.Backend{}
	regular := []*config.Backend{}

	for _, b := range cfg.Backends {
		if isShadowBackend(b) {
			shadow = append(shadow, b)
			continue
		}
		regular = append(regular, b)
	}

	cfg.Backends = regular

	p, err = s.f.New(cfg)

	if len(shadow) > 0 {
		cfg.Backends = shadow
		pShadow, _ := s.f.New(cfg)
		p = ShadowMiddleware(p, pShadow)
	}

	return
}

// NewShadowFactory 使用提供的工厂创建一个新的shadowFactory
func NewShadowFactory(f Factory) Factory {
	return shadowFactory{f}
}

// ShadowMiddleware 是一个创建shadowProxy的中间件
func ShadowMiddleware(next ...Proxy) Proxy {
	switch len(next) {
	case 0:
		panic(ErrNotEnoughProxies)
	case 1:
		return next[0]
	case 2:
		return NewShadowProxy(next[0], next[1])
	default:
		panic(ErrTooManyProxies)
	}
}

// NewShadowProxy 返回一个向p1和p2发送请求但忽略p2响应的代理
func NewShadowProxy(p1, p2 Proxy) Proxy {
	return func(ctx context.Context, request *Request) (*Response, error) {
		go p2(newcontextWrapper(ctx), CloneRequest(request))
		return p1(ctx, request)
	}
}

func isShadowBackend(c *config.Backend) bool {
	if v, ok := c.ExtraConfig[Namespace]; ok {
		if e, ok := v.(map[string]interface{}); ok {
			if v, ok := e[shadowKey]; ok {
				c, ok := v.(bool)
				return ok && c
			}
		}
	}
	return false
}

type contextWrapper struct {
	context.Context
	data context.Context
}

func (c contextWrapper) Value(key interface{}) interface{} {
	return c.data.Value(key)
}

func newcontextWrapper(data context.Context) contextWrapper {
	return contextWrapper{
		Context: context.Background(),
		data:    data,
	}
}
