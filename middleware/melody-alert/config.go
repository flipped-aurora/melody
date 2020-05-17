package alert

import (
	"errors"
	"melody/config"
)

const (
	namespace = "melody_alert"
)

func NewAPIChecker(cfg *config.EndpointConfig) (Checker, error) {
	m, err := parseConfig(cfg.ExtraConfig)
	if err != nil {
		return nil, err
	}
	m["api"] = cfg.Endpoint
	checker, err := newChecker(m)
	if err != nil {
		return nil, err
	}
	return checker, nil
}

func NewBootChecker(cfg config.ExtraConfig) (Checker, error) {
	m, err := parseConfig(cfg)
	if err != nil {
		return nil, err
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
