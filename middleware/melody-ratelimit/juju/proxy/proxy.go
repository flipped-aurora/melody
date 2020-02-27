package proxy

import (
	"context"

	"melody/config"
	"melody/proxy"

	melodyrate "melody/middleware/melody-ratelimit"
	"melody/middleware/melody-ratelimit/juju"
)

// Namespace is the key to use to store and access the custom config data for the proxy
const Namespace = "melody_ratelimit_proxy"

// Config is the custom config struct containing the params for the limiter
type Config struct {
	MaxRate  float64
	Capacity int64
}

// BackendFactory adds a ratelimiting middleware wrapping the internal factory
func BackendFactory(next proxy.BackendFactory) proxy.BackendFactory {
	return func(cfg *config.Backend) proxy.Proxy {
		return NewMiddleware(cfg)(next(cfg))
	}
}

// NewMiddleware builds a middleware based on the extra config params or fallbacks to the next proxy
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

// ZeroCfg is the zero value for the Config struct
var ZeroCfg = Config{}

// ConfigGetter implements the config.ConfigGetter interface. It parses the extra config for the
// rate adapter and returns a ZeroCfg if something goes wrong.
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
