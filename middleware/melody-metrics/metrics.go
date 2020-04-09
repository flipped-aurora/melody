package metrics

import (
	"context"
	"melody/config"
	"melody/logging"
	"time"

	"github.com/rcrowley/go-metrics"
)

const (
	// Namespace 标识是否开启metrics
	Namespace = "melody_metrics"
	// DefaultListenAddr metrics服务默认监听的端口
	DefaultListenAddr = ":8090"
)

var (
	percentiles   = []float64{0.1, 0.25, 0.5, 0.75, 0.9, 0.95, 0.99}
	defaultSample = func() metrics.Sample { return metrics.NewUniformSample(1028) }
)

// Config metrics 服务默认配置
type Config struct {
	ProxyDisabled    bool
	RouterDisabled   bool
	BackendDisabled  bool
	CollectionTime   time.Duration
	ListenAddr       string
	EndpointDisabled bool
}

// Metrics metrics结构对象
type Metrics struct {
	// 为Config提供的计数器
	Config *Config
	// 为Proxy提供的计数器
	Proxy *ProxyMetrics
	// 为Router模块提供的计数器
	Router *RouterMetrics
	// 注册计时器
	Registry       *metrics.Registry
	latestSnapshot Stats
}

// New 返回一个metrics的运行实例
func New(ctx context.Context, e config.ExtraConfig, logger logging.Logger) *Metrics {
	registry := metrics.NewPrefixedRegistry("melody.")

	var metricsConfig *Config
	if c, ok := GetConfig(e).(*Config); ok {
		metricsConfig = c
	}

	if metricsConfig == nil {
		registry = NewNullRegistry()
		return &Metrics{
			Proxy:    &ProxyMetrics{},
			Router:   &RouterMetrics{},
			Registry: &registry,
		}
	}

	m := Metrics{
		Config:         metricsConfig,
		Proxy:          NewProxyMetrics(&registry),  // melody.proxy.
		Router:         NewRouterMetrics(&registry), // melody.router.
		Registry:       &registry,
		latestSnapshot: NewStats(),
	}

	m.processMetrics(ctx, m.Config.CollectionTime, logger) // melody.service.

	return &m
}

// debugGC runtimeMem
func (m *Metrics) processMetrics(ctx context.Context, duration time.Duration, logger logging.Logger) {
	r := metrics.NewPrefixedChildRegistry(*(m.Registry), "service.")

	metrics.RegisterDebugGCStats(r)
	metrics.RegisterRuntimeMemStats(r)

	go func() {
		ticket := time.NewTicker(duration)

		for {
			select {
			case <-ticket.C:
				metrics.CaptureDebugGCStatsOnce(r)
				metrics.CaptureRuntimeMemStatsOnce(r)
				m.Router.Aggregate() // 统计 router 的连接情况
				m.latestSnapshot = m.TakeSnapshot()
			case <-ctx.Done():
				return
			}
		}
	}()
}

// TakeSnapshot 制作一个上下文瞬间
func (m *Metrics) TakeSnapshot() Stats {
	sta := NewStats()

	(*m.Registry).Each(func(s string, i interface{}) {
		switch metric := i.(type) {
		case metrics.Counter:
			sta.Counters[s] = metric.Count()
		case metrics.Gauge:
			sta.Gauges[s] = metric.Value()
		case metrics.Histogram:
			sta.Histograms[s] = HistogramData{
				Max:         metric.Max(),
				Min:         metric.Min(),
				Mean:        metric.Mean(),
				Stddev:      metric.StdDev(),
				Variance:    metric.Variance(),
				Percentiles: metric.Percentiles(percentiles),
			}
			metric.Clear()
		}
	})

	return sta
}

func (m *Metrics) GetSnapshot() Stats {
	return m.latestSnapshot
}

// NewNullRegistry 返回一个Null registry
func NewNullRegistry() metrics.Registry {
	return &NullRegistry{}
}

// GetConfig 从ExtraConfig中提取metrics的配置
func GetConfig(e config.ExtraConfig) interface{} {
	v, ok := e[Namespace]
	if !ok {
		return nil
	}

	temp, ok := v.(map[string]interface{})
	if !ok {
		return nil
	}

	config := new(Config)
	config.CollectionTime = time.Minute
	if t, ok := temp["collection_time"]; ok {
		if d, err := time.ParseDuration(t.(string)); err == nil {
			config.CollectionTime = d
		}
	}

	config.ListenAddr = DefaultListenAddr
	if a, ok := temp["listen_address"]; ok {
		if ad, ok := a.(string); ok {
			config.ListenAddr = ad
		}
	}

	config.ProxyDisabled = getBool(temp, "proxy_disabled")
	config.RouterDisabled = getBool(temp, "router_disabled")
	config.BackendDisabled = getBool(temp, "backend_disabled")
	config.EndpointDisabled = getBool(temp, "endpoint_disabled")

	return config
}

func getBool(temp map[string]interface{}, s string) bool {
	if t, ok := temp[s]; ok {
		if v, ok := t.(bool); ok {
			return v
		}
	}
	return false
}

type NullRegistry struct{}

func (n *NullRegistry) Each(func(string, interface{})) {}

func (n *NullRegistry) Get(string) interface{} {
	return nil
}

func (n *NullRegistry) GetAll() map[string]map[string]interface{} {
	return nil
}

func (n *NullRegistry) GetOrRegister(string, interface{}) interface{} {
	return nil
}

func (n *NullRegistry) Register(string, interface{}) error {
	return nil
}

func (n *NullRegistry) RunHealthchecks() {}

func (n *NullRegistry) Unregister(string) {}

func (n *NullRegistry) UnregisterAll() {}
