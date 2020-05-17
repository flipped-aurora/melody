package alert

import (
	"errors"
	"melody/config"
)

const (
	namespace = "melody_alert"
)

func NewChecker(cfg *config.ServiceConfig) (Checker, error) {
	// 解析Service
	m, err := parseConfig(cfg.ExtraConfig)
	if err != nil {
		return nil, err
	}

	// 解析Endpoint
	for _, endpointConfig := range cfg.Endpoints {
		endpointM, err := parseConfig(endpointConfig.ExtraConfig)
		if err == nil {
			m[endpointConfig.Endpoint] = endpointM
		}
	}

	checker, err := newChecker(m)
	if err != nil {
		return nil, err
	}
	return checker, nil
}

func parseConfig(extraConfig config.ExtraConfig) (map[string]interface{}, error) {
	if _, ok := extraConfig[namespace]; !ok {
		return nil, errors.New("no melody_alert")
	}

	if fm, ok := extraConfig[namespace].(map[string]interface{}); !ok {
		return nil, errors.New("no fields")
	} else {
		return fm, nil
	}
}
