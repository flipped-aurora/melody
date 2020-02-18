package metrics

import "github.com/rcrowley/go-metrics"

type RouterMetrics struct {
	ProxyMetrics
	connected         metrics.Counter
	disconnected      metrics.Counter
	connectedTotal    metrics.Counter
	disconnectedTotal metrics.Counter
	connectedGauge    metrics.Gauge
	disconnectedGauge metrics.Gauge
}

func NewRouterMetrics(parent *metrics.Registry) *RouterMetrics {
	r := metrics.NewPrefixedChildRegistry(*parent, "router.")

	return &RouterMetrics{
		ProxyMetrics:      ProxyMetrics{register: r},
		connected:         metrics.NewRegisteredCounter("connected", r),
		disconnected:      metrics.NewRegisteredCounter("disconnected", r),
		connectedTotal:    metrics.NewRegisteredCounter("connected-total", r),
		disconnectedTotal: metrics.NewRegisteredCounter("disconnected-total", r),
		connectedGauge:    metrics.NewRegisteredGauge("connected-gauge", r),
		disconnectedGauge: metrics.NewRegisteredGauge("disconnected-gauge", r),
	}
}

func (rm *RouterMetrics) Aggregate() {
	con := rm.connected.Count()
	rm.connectedGauge.Update(con)
	rm.connectedTotal.Inc(con)
	rm.connected.Clear()
	discon := rm.disconnected.Count()
	rm.disconnectedGauge.Update(discon)
	rm.disconnectedTotal.Inc(discon)
	rm.disconnected.Clear()
}
