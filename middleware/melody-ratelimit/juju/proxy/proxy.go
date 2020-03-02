package proxy

import (
	"context"

	"melody/config"
	"melody/proxy"

	melodyrate "melody/middleware/melody-ratelimit"
	"melody/middleware/melody-ratelimit/juju"
)

// Namespace 命名空间
const Namespace = "melody_ratelimit_proxy"

// Config 是包含限定符参数的自定义配置结构
type Config struct {
	MaxRate  float64
	Capacity int64
}

// BackendFactory 添加了一个包装内部工厂的速率限制中间件
func BackendFactory(next proxy.BackendFactory) proxy.BackendFactory {
	return func(cfg *config.Backend) proxy.Proxy {
		return NewMiddleware(cfg)(next(cfg))
	}
}

// NewMiddleware 基于对下一个代理的额外配置参数或回退构建中间件
func NewMiddleware(remote *config.Backend) proxy.Middleware {
	cfg := ConfigGetter(remote.ExtraConfig).(Config)
	if cfg == ZeroCfg || cfg.MaxRate <= 0 {
		return proxy.EmptyMiddleware
	}
	tb := juju.NewLimiter(cfg.MaxRate, cfg.Capacity)
	return func(next ...proxy.Proxy) proxy.Proxy {
		if len(next) > 1 {
			panic(proxy.ErrTooManyProxies)
		}
		return func(ctx context.Context, request *proxy.Request) (*proxy.Response, error) {
			if !tb.Allow() {
				return nil, melodyrate.ErrLimited
			}
			return next[0](ctx, request)
		}
	}
}

// ZeroCfg 是配置结构的空值
var ZeroCfg = Config{}

//ConfigGetter 实现config.ConfigGetter接口。它解析速率适配器的额外配置，如果出了问题，则返回一个ZeroCfg。
func ConfigGetter(e config.ExtraConfig) interface{} {
	v, ok := e[Namespace]
	if !ok {
		return ZeroCfg
	}
	tmp, ok := v.(map[string]interface{})
	if !ok {
		return ZeroCfg
	}
	cfg := Config{}
	if v, ok := tmp["maxRate"]; ok {
		switch val := v.(type) {
		case float64:
			cfg.MaxRate = val
		case int:
			cfg.MaxRate = float64(val)
		case int64:
			cfg.MaxRate = float64(val)
		}
	}
	if v, ok := tmp["capacity"]; ok {
		switch val := v.(type) {
		case int64:
			cfg.Capacity = val
		case int:
			cfg.Capacity = int64(val)
		case float64:
			cfg.Capacity = int64(val)
		}
	}
	return cfg
}
