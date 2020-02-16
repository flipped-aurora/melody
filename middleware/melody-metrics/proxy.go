package metrics

import "github.com/rcrowley/go-metrics"

type ProxyMetrics struct {
	register metrics.Registry
}

func NewProxyMetrics(parent *metrics.Registry) *ProxyMetrics {
	m := metrics.NewPrefixedChildRegistry(*parent, "proxy.")
	return &ProxyMetrics{register:m}
}
