package gobreaker

import (
	"fmt"
	"github.com/sony/gobreaker"
	"melody/config"
	"melody/logging"
	"time"
)

const Namespace = "melody_circuitbreaker"

type Config struct {
	//给定的时间间隔（秒）
	Interval int
	//等待时间窗口（秒）
	Timeout int
	//连续故障数
	MaxErrors int
	//断路器状态发生改变时，是否log
	LogStatusChange bool
}

// 空实现
var DefaultCfg = Config{}

// 获取断路器配置
func ConfigGetter(e config.ExtraConfig) interface{} {
	v, ok := e[Namespace]
	if !ok {
		return DefaultCfg
	}
	temp, ok := v.(map[string]interface{})
	if !ok {
		return DefaultCfg
	}
	cfg := Config{}
	if i, ok := temp["interval"]; ok {
		switch in := i.(type) {
		case int:
			cfg.Interval = in
		case int64:
			cfg.Interval = int(in)
		}
	}
	if v, ok := temp["timeout"]; ok {
		switch i := v.(type) {
		case int:
			cfg.Timeout = i
		case float64:
			cfg.Timeout = int(i)
		}
	}
	if v, ok := temp["maxErrors"]; ok {
		switch i := v.(type) {
		case int:
			cfg.MaxErrors = i
		case float64:
			cfg.MaxErrors = int(i)
		}
	}
	value, ok := temp["logStatusChange"].(bool)
	cfg.LogStatusChange = ok && value

	return cfg
}

// 基于sony的断路器，包装断路器对象
func NewCircuitBreaker(config Config, logger logging.Logger) *gobreaker.CircuitBreaker {
	settings := gobreaker.Settings{
		Name:     "Melody CircuitBreaker",
		Interval: time.Duration(config.Interval) * time.Second,
		Timeout:  time.Duration(config.Timeout) * time.Second,
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			return counts.ConsecutiveFailures > uint32(config.MaxErrors)
		},
	}

	if config.LogStatusChange {
		settings.OnStateChange = func(name string, from gobreaker.State, to gobreaker.State) {
			logger.Warning(fmt.Sprintf("circuit breaker named '%s' went from '%s' to '%s'", name, from.String(), to.String()))
		}
	}

	return gobreaker.NewCircuitBreaker(settings)
}
