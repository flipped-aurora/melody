package metrics

import (
	"context"
	"github.com/rcrowley/go-metrics"
	"melody/config"
	"melody/proxy"
	"strconv"
	"strings"
	"time"
)

type ProxyMetrics struct {
	register metrics.Registry
}

func (m *Metrics) NewProxyFactory(segmentName string, next proxy.Factory) proxy.FactoryFunc {
	if m.Config == nil || m.Config.ProxyDisabled {
		return next.New
	}
	return func(cfg *config.EndpointConfig) (proxy.Proxy, error) {
		next, err := next.New(cfg)
		if err != nil {
			return proxy.NoopProxy, err
		}
		return m.NewProxyMiddleware(segmentName, cfg.Endpoint)(next), nil
	}
}

func (m *Metrics) NewBackendFactory(prefixName string, next proxy.BackendFactory) proxy.BackendFactory {
	if m.Config == nil || m.Config.BackendDisabled {
		return next
	}
	return func(backend *config.Backend) proxy.Proxy {
		return m.NewProxyMiddleware(prefixName, backend.URLPattern)(next(backend))
	}
}

func (m *Metrics) NewProxyMiddleware(layer, name string) proxy.Middleware {
	return NewProxyMiddleware(layer, name, m.Proxy)
}

func NewProxyMetrics(parent *metrics.Registry) *ProxyMetrics {
	m := metrics.NewPrefixedChildRegistry(*parent, "proxy.")
	return &ProxyMetrics{register: m}
}

func NewProxyMiddleware(layer, name string, pm *ProxyMetrics) proxy.Middleware {
	registerProxyMiddlewareMetrics(layer, name, pm)
	return func(p ...proxy.Proxy) proxy.Proxy {
		if len(p) > 1 {
			panic(proxy.ErrTooManyProxies)
		}

		return func(ctx context.Context, request *proxy.Request) (response *proxy.Response, err error) {
			begin := time.Now()
			resp, err := p[0](ctx, request)

			go func(duration int64, resp *proxy.Response, err error) {
				errored := strconv.FormatBool(err != nil)
				complete := strconv.FormatBool(resp != nil && resp.IsComplete)
				labels := "layer." + layer + ".name." + name + ".complete." + complete + ".error." + errored
				pm.Counter("requests." + labels).Inc(1)
				pm.Histogram("latency." + labels).Update(duration)
			}(time.Since(begin).Nanoseconds(), resp, err)

			return resp, err
		}
	}
}

func registerProxyMiddlewareMetrics(layer, name string, pm *ProxyMetrics) {
	labels := "layer." + layer + ".name." + name
	for _, complete := range []string{"true", "false"} {
		for _, errored := range []string{"true", "false"} {
			metrics.GetOrRegisterCounter("requests."+labels+".complete."+complete+".error."+errored, pm.register)

			metrics.GetOrRegisterHistogram("latency."+labels+".complete."+complete+".error."+errored, pm.register, defaultSample())
		}
	}
}

// 获取或注册直方图
func (rm *ProxyMetrics) Histogram(labels ...string) metrics.Histogram {
	return metrics.GetOrRegisterHistogram(strings.Join(labels, "."), rm.register, defaultSample())
}

// 获取或注册计数器
func (rm *ProxyMetrics) Counter(labels ...string) metrics.Counter {
	return metrics.GetOrRegisterCounter(strings.Join(labels, "."), rm.register)
}

