package logstash

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"melody/config"
	"melody/logging"
	gologging "melody/middleware/melody-gologging"
	"os"
	"time"
)

const (
	Namespace = "melody_logstash"
	Prefix    = "prefix"

	LEVEL_DEBUG    = "DEBUG"
	LEVEL_INFO     = "INFO"
	LEVEL_WARNING  = "WARNING"
	LEVEL_ERROR    = "ERROR"
	LEVEL_CRITICAL = "CRITICAL"

	ISO_8601 = "2006-01-02T15:04:05.000000-07:00"
)

var (
	now               = time.Now
	hostname          = "localhost"
	loggingPattern    = "%{message}"
	ErrorConfigError  = errors.New("not found extra config about melody-logstash module")
	ErrorNothingToLog = errors.New("there are nothing to log")
)

type LogLevel string

func init() {
	name, err := os.Hostname()
	if err != nil {
		hostname = name
	}
}

func NewLogger(extra config.ExtraConfig, ws ...io.Writer) (logging.Logger, error) {
	_, ok := extra[Namespace]
	if !ok {
		return nil, ErrorConfigError
	}

	serviceName := "MELODY"
	gologging.DefaultPattern = loggingPattern
	if log, ok := extra[gologging.Namespace]; ok {
		if section, ok := log.(map[string]interface{}); ok {
			if pre, ok := section[Prefix]; ok {
				if str, ok := pre.(string); ok {
					serviceName = str
					delete(section, Prefix)
				}
			}
		}
	}
	logger, err := gologging.NewLogger(extra, ws...)
	if err != nil {
		return nil, err
	}
	return &Logger{
		logger:      logger,
		serviceName: serviceName,
	}, nil
}

// Logger is a wrapper over a github.com/devopsfaith/krakend/logging logger
type Logger struct {
	logger      logging.Logger
	serviceName string
}

// Debug implements the logger interface
func (l *Logger) Debug(v ...interface{}) {
	data, err := l.format(LEVEL_DEBUG, v...)
	if err != nil {
		return
	}
	l.logger.Debug(string(data))
}

// Info implements the logger interface
func (l *Logger) Info(v ...interface{}) {
	data, err := l.format(LEVEL_INFO, v...)
	if err != nil {
		return
	}
	l.logger.Info(string(data))
}

// Warning implements the logger interface
func (l *Logger) Warning(v ...interface{}) {
	data, err := l.format(LEVEL_WARNING, v...)
	if err != nil {
		return
	}
	l.logger.Warning(string(data))
}

// Error implements the logger interface
func (l *Logger) Error(v ...interface{}) {
	data, err := l.format(LEVEL_ERROR, v...)
	if err != nil {
		return
	}
	l.logger.Error(string(data))
}

// Critical implements the logger interface
func (l *Logger) Critical(v ...interface{}) {
	data, err := l.format(LEVEL_CRITICAL, v...)
	if err != nil {
		return
	}
	l.logger.Critical(string(data))
}

// Fatal implements the logger interface
func (l *Logger) Fatal(v ...interface{}) {
	data, err := l.format(LEVEL_CRITICAL, v...)
	if err != nil {
		return
	}
	l.logger.Fatal(string(data))
}

func (l *Logger) format(logLevel LogLevel, v ...interface{}) ([]byte, error) {
	if len(v) == 0 {
		return []byte{}, ErrorNothingToLog
	}

	msg, ok := v[0].(string)
	if !ok {
		msg = fmt.Sprintf("%+v", v[0])
	}
	record := map[string]interface{}{}
	if len(v) > 1 {
		for _, ctx := range v[1:] {
			switch value := ctx.(type) {
			case map[string]interface{}:
				for k, item := range value {
					record[k] = item
				}
			case int:
				msg = fmt.Sprintf("%s %d", msg, value)
			case bool:
				msg = fmt.Sprintf("%s %t", msg, value)
			case float64:
				msg = fmt.Sprintf("%s %f", msg, value)
			case string:
				msg += " " + value
			//case *proxy.Request:
			//	record["proxy.Request"] = value
			//case *proxy.Response:
			//	record["proxy.Response"] = value
			default:
				record[fmt.Sprintf("%T", ctx)] = ctx
			}
		}
	}

	record["@version"] = 1
	record["@timestamp"] = now().Format(ISO_8601)
	record["module"] = l.serviceName
	record["host"] = hostname
	record["message"] = msg
	record["level"] = logLevel

	return json.Marshal(record)
}
