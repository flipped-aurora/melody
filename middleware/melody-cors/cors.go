package cors

import (
	"melody/config"
	"time"
)

// Namespace the key of corss domain
const Namespace = "melody_cors"

// Config holds the configuration of CORS
type Config struct {
	AllowOrigins     []string
	AllowMethods     []string
	AllowHeaders     []string
	ExposeHeaders    []string
	AllowCredentials bool
	MaxAge           time.Duration
}

// GetConfig 获取cors配置strcut
func GetConfig(e config.ExtraConfig) interface{} {
	v, ok := e[Namespace]
	if !ok {
		return nil
	}

	temp, ok := v.(map[string]interface{})
	if !ok {
		return nil
	}

	cfg := Config{}
	cfg.AllowOrigins = getList(temp, "allow_origins")
	cfg.AllowMethods = getList(temp, "allow_methods")
	cfg.AllowHeaders = getList(temp, "allow_headers")
	cfg.ExposeHeaders = getList(temp, "expose_headers")
	if cr, ok := temp["allow_credentials"]; ok {
		if cre, ok := cr.(bool); ok {
			cfg.AllowCredentials = cre
		}
	}

	if a, ok := temp["max_age"]; ok {
		if d, err := time.ParseDuration(a.(string)); err == nil {
			cfg.MaxAge = d
		}
	}

	return cfg
}

// getList return array data from map via key
func getList(data map[string]interface{}, name string) []string {
	out := []string{}
	if vs, ok := data[name]; ok {
		if v, ok := vs.([]interface{}); ok {
			for _, s := range v {
				if j, ok := s.(string); ok {
					out = append(out, j)
				}
			}
		}
	}
	return out
}
