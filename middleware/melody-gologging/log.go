package gologging

import (
	"fmt"
	"io"
	"melody/config"
	"melody/logging"
	"os"

	oplogging "github.com/op/go-logging"
)

const (
	Namespace    = "melody_gologging"
	Level        = "level"
	Syslog       = "syslog"
	StdOut       = "stdout"
	Prefix       = "prefix"
	Format       = "format"
	CustomFormat = "custom_format"
)

var (
	DefaultPattern           = `%{time:2006/01/02 - 15:04:05.000} %{color}▶ %{level:.6s}%{color:reset} %{message}`
	LogstashPattern          = `{"@timestamp":"%{time:200-01-02T15:04:05.000+00:00}", "@version": 1, "level": "%{level}", "message": "%{message}", "module": "%{module}"}`
	ActivePattren            = DefaultPattern
	ErrorWrongConfig         = fmt.Errorf("not found extra config about melody-gologging module")
	defaultFormatterSelector = func(io.Writer) string { return ActivePattren }
)

//Config contains config of logger designed by user
type Config struct {
	Level        string
	Syslog       bool
	StdOut       bool
	Prefix       string
	Format       string
	CustomFormat string
}

type Logger struct {
	logger *oplogging.Logger
}

func (l Logger) Debug(v ...interface{}) {
	l.logger.Debug(v)
}

func (l Logger) Info(v ...interface{}) {
	l.logger.Info(v)
}

func (l Logger) Warning(v ...interface{}) {
	l.logger.Warning(v)
}

func (l Logger) Error(v ...interface{}) {
	l.logger.Error(v)
}

func (l Logger) Critical(v ...interface{}) {
	l.logger.Critical(v)
}

func (l Logger) Fatal(v ...interface{}) {
	l.logger.Fatal(v)
}

func NewLogger(config config.ExtraConfig, ws ...io.Writer) (logging.Logger, error) {
	cfg, ok := GetConfig(config).(Config)
	if !ok {
		return nil, ErrorWrongConfig
	}

	module := "MELODY"
	logger := oplogging.MustGetLogger(module)

	if cfg.StdOut {
		ws = append(ws, os.Stdout)
	}

	//if cfg.Syslog {
	//	w, err := syslog.New(syslog.LOG_CRIT, cfg.Prefix)
	//	if err != nil {
	//		return nil, err
	//	}
	//
	//	ws = append(ws, w)
	//}

	switch cfg.Format {
	case "logstash":
		ActivePattren = LogstashPattern
		cfg.Prefix = ""
	case "custom":
		ActivePattren = cfg.CustomFormat
		cfg.Prefix = ""
	}

	//We need many Backend because of many writer
	var backends []oplogging.Backend
	for _, w := range ws {
		backend := oplogging.NewLogBackend(w, cfg.Prefix, 0)
		pattern := defaultFormatterSelector(w)
		format := oplogging.MustStringFormatter(pattern)
		backendLeveled := oplogging.AddModuleLevel(oplogging.NewBackendFormatter(backend, format))
		level, err := oplogging.LogLevel(cfg.Level)
		if err != nil {
			return nil, err
		}

		backendLeveled.SetLevel(level, module)
		backends = append(backends, backendLeveled)
	}

	oplogging.SetBackend(backends...)

	return Logger{logger: logger}, nil
}

//GetConfig put extra config into config struct
func GetConfig(extraConfig config.ExtraConfig) interface{} {
	v, ok := extraConfig[Namespace]
	if !ok {
		return nil
	}

	m, ok := v.(map[string]interface{})
	if !ok {
		return nil
	}

	cfg := Config{}
	if v, ok := m[StdOut]; ok {
		cfg.StdOut = v.(bool)
	}
	if v, ok := m[Syslog]; ok {
		cfg.Syslog = v.(bool)
	}
	if v, ok := m[Level]; ok {
		cfg.Level = v.(string)
	}
	if v, ok := m[Prefix]; ok {
		cfg.Prefix = v.(string)
	}
	if v, ok := m[Format]; ok {
		cfg.Format = v.(string)
	}
	if v, ok := m[CustomFormat]; ok {
		cfg.CustomFormat = v.(string)
	}

	return cfg
}

func UpdateFormatSelector(f func(io.Writer) string) {
	defaultFormatterSelector = f
}
