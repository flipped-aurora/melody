package opencensus

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"sync"
	"time"

	"melody/config"

	"go.opencensus.io/plugin/ochttp"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/trace"
)

type ExporterFactory func(context.Context, Config) (interface{}, error)

func RegisterExporterFactories(ef ExporterFactory) {
	mu.Lock()
	exporterFactories = append(exporterFactories, ef)
	mu.Unlock()
}

func Register(ctx context.Context, srvCfg config.ServiceConfig, vs ...*view.View) error {
	cfg, err := parseCfg(srvCfg)
	if err != nil {
		return err
	}

	err = errSingletonExporterFactoriesRegister
	registerOnce.Do(func() {
		register.ExporterFactories(ctx, *cfg, exporterFactories)

		err = register.Register(ctx, *cfg, vs)
		if err != nil {
			return
		}

		if cfg.EnabledLayers != nil {
			enabledLayers = *cfg.EnabledLayers
			return
		}

		enabledLayers = EnabledLayers{true, true, true}
	})

	return err
}

type composableRegister struct {
	viewExporter       func(exporters ...view.Exporter)
	traceExporter      func(exporters ...trace.Exporter)
	registerViews      func(views ...*view.View) error
	setDefaultSampler  func(rate int)
	setReportingPeriod func(d time.Duration)
}

func (c *composableRegister) ExporterFactories(ctx context.Context, cfg Config, fs []ExporterFactory) {
	viewExporters := []view.Exporter{}
	traceExporters := []trace.Exporter{}

	for _, f := range fs {
		e, err := f(ctx, cfg)
		if err != nil {
			continue
		}
		if ve, ok := e.(view.Exporter); ok {
			viewExporters = append(viewExporters, ve)
		}
		if te, ok := e.(trace.Exporter); ok {
			traceExporters = append(traceExporters, te)
		}
	}

	c.viewExporter(viewExporters...)
	c.traceExporter(traceExporters...)
}

func (c composableRegister) Register(ctx context.Context, cfg Config, vs []*view.View) error {
	if len(vs) == 0 {
		vs = DefaultViews
	}

	c.setDefaultSampler(cfg.SampleRate)
	c.setReportingPeriod(time.Duration(cfg.ReportingPeriod) * time.Second)

	return c.registerViews(vs...)
}

type Config struct {
	SampleRate      int            `json:"sample_rate"`
	ReportingPeriod int            `json:"reporting_period"`
	EnabledLayers   *EnabledLayers `json:"enabled_layers"`
	Exporters       struct {
		InfluxDB *struct {
			Address      string `json:"address"`
			Username     string `json:"username"`
			Password     string `json:"password"`
			Timeout      string `json:"timeout"`
			PingEnabled  bool   `json:"ping"`
			Database     string `json:"db"`
			InstanceName string `json:"service_name"`
			BufferSize   int    `json:"buffer_size"`
		} `json:"influxdb"`
		Zipkin *struct {
			CollectorURL string `json:"collector_url"`
			ServiceName  string `json:"service_name"`
			IP           string `json:"ip"`
			Port         int    `json:"port"`
		} `json:"zipkin"`
		Jaeger *struct {
			Endpoint    string `json:"endpoint"`
			ServiceName string `json:"service_name"`
		} `json:"jaeger"`
		Prometheus *struct {
			Port int `json:"port"`
		} `json:"prometheus"`
		Logger *struct{} `json:"logger"`
		Xray   *struct {
			UseEnv    bool   `json:"use_env"`
			Region    string `json:"region"`
			AccessKey string `json:"access_key_id"`
			SecretKey string `json:"secret_access_key"`
			Version   string `json:"version"`
		} `json:"xray"`
		Stackdriver *struct {
			ProjectID     string            `json:"project_id"`
			MetricPrefix  string            `json:"metric_prefix"`
			DefaultLabels map[string]string `json:"default_labels"`
		} `json:"stackdriver"`
	} `json:"exporters"`
}

const (
	ContextKey = "opencensus-request-span"
	Namespace  = "melody_opencensus"
)

var (
	DefaultViews = []*view.View{
		ochttp.ClientSentBytesDistribution,
		ochttp.ClientReceivedBytesDistribution,
		ochttp.ClientRoundtripLatencyDistribution,
		ochttp.ClientCompletedCount,

		ochttp.ServerRequestCountView,
		ochttp.ServerRequestBytesView,
		ochttp.ServerResponseBytesView,
		ochttp.ServerLatencyView,
		ochttp.ServerRequestCountByMethod,
		ochttp.ServerResponseCountByStatusCode,
	}

	exporterFactories                     = []ExporterFactory{}
	errNoExtraConfig                      = errors.New("no extra config defined for the opencensus module")
	errSingletonExporterFactoriesRegister = errors.New("expecting only one exporter factory registration per instance")
	mu                                    = new(sync.RWMutex)
	register                              = composableRegister{
		viewExporter:       registerViewExporter,
		traceExporter:      registerTraceExporter,
		setDefaultSampler:  setDefaultSampler,
		setReportingPeriod: setReportingPeriod,
		registerViews:      registerViews,
	}
	registerOnce  = new(sync.Once)
	enabledLayers EnabledLayers
)

type EnabledLayers struct {
	Router  bool `json:"router"`
	Pipe    bool `json:"pipe"`
	Backend bool `json:"backend"`
}

func IsRouterEnabled() bool {
	return enabledLayers.Router
}

func IsPipeEnabled() bool {
	return enabledLayers.Pipe
}

func IsBackendEnabled() bool {
	return enabledLayers.Backend
}

func parseCfg(srvCfg config.ServiceConfig) (*Config, error) {
	cfg := new(Config)
	tmp, ok := srvCfg.ExtraConfig[Namespace]
	if !ok {
		return nil, errNoExtraConfig
	}
	buf := new(bytes.Buffer)
	json.NewEncoder(buf).Encode(tmp)
	if err := json.NewDecoder(buf).Decode(cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}

func fromContext(ctx context.Context) *trace.Span {
	span := trace.FromContext(ctx)
	if span == nil {
		span, _ = ctx.Value(ContextKey).(*trace.Span)
	}
	return span
}

func registerViewExporter(exporters ...view.Exporter) {
	for _, e := range exporters {
		view.RegisterExporter(e)
	}
}

func registerTraceExporter(exporters ...trace.Exporter) {
	for _, e := range exporters {
		trace.RegisterExporter(e)
	}
}

func setDefaultSampler(rate int) {
	var sampler trace.Sampler
	switch {
	case rate <= 0:
		sampler = trace.NeverSample()
	case rate >= 100:
		sampler = trace.AlwaysSample()
	default:
		sampler = trace.ProbabilitySampler(float64(rate) / 100.0)
	}
	trace.ApplyConfig(trace.Config{DefaultSampler: sampler})
}

func setReportingPeriod(d time.Duration) {
	view.SetReportingPeriod(d)
}

func registerViews(views ...*view.View) error {
	return view.Register(views...)
}
