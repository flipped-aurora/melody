package consul

import (
	"fmt"
	"melody/config"
)

type Config struct {
	Address, Name string
	Port          int
	Tags          []string
}

// Namespace is the key to use to store and access the custom config data
const Namespace = "github_com/letgoapp/krakend-consul"

var (
	// ErrNoConfig is the error to be returned when there is no config with the consul namespace
	ErrNoConfig = fmt.Errorf("unable to create the consul client: no config")
	// ErrBadConfig is the error to be returned when the config is not well defined
	ErrBadConfig = fmt.Errorf("unable to create the consul client with the received config")
	// ErrNoMachines is the error to be returned when the config has not defined one or more servers
	ErrNoMachines = fmt.Errorf("unable to create the consul client without a set of servers")
)

func parse(e config.ExtraConfig, port int) (Config, error) {
	cfg := Config{
		Name: "krakend",
		Port: port,
	}
	v, ok := e[Namespace]
	if !ok {
		return cfg, ErrNoConfig
	}
	tmp, ok := v.(map[string]interface{})
	if !ok {
		return cfg, ErrBadConfig
	}
	a, ok := tmp["address"]
	if !ok {
		return cfg, ErrNoMachines
	}
	cfg.Address, ok = a.(string)

	if !ok {
		return cfg, ErrNoMachines
	}

	if a, ok = tmp["name"]; ok {
		cfg.Name, ok = a.(string)
	}

	cfg.Tags = parseTags(tmp)

	return cfg, nil
}

func parseTags(cfg map[string]interface{}) []string {
	result := []string{}
	tags, ok := cfg["tags"]
	if !ok {
		return result
	}
	tgs, ok := tags.([]interface{})
	if !ok {
		return result
	}
	for _, tg := range tgs {
		if t, ok := tg.(string); ok {
			result = append(result, t)
		}
	}

	return result
}
