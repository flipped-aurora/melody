package router

import (
	"fmt"

	"melody/config"
)

// Namespace 命名空间是用来存储和访问路由器自定义配置数据
const Namespace = "melody_ratelimit_router"

// Config 是包含路由器中间件参数的自定义配置结构
type Config struct {
	MaxRate       int64
	Strategy      string
	ClientMaxRate int64
	Key           string
}

// ZeroCfg ZeroCfg 是Config struct的零值
var ZeroCfg = Config{}

// ConfigGetter 实现config.ConfigGetter接口。它解析速率适配器的extra config，如果出了问题，则返回一个ZeroCfg。
func ConfigGetter(e config.ExtraConfig) interface{} {
	v, ok := e[Namespace]
	if !ok {
		return ZeroCfg
	}
	tmp, ok := v.(map[string]interface{})
	if !ok {
		return ZeroCfg
	}
	cfg := Config{}
	if v, ok := tmp["maxRate"]; ok {
		switch val := v.(type) {
		case int64:
			cfg.MaxRate = val
		case int:
			cfg.MaxRate = int64(val)
		case float64:
			cfg.MaxRate = int64(val)
		}
	}
	if v, ok := tmp["strategy"]; ok {
		cfg.Strategy = fmt.Sprintf("%v", v)
	}
	if v, ok := tmp["clientMaxRate"]; ok {
		switch val := v.(type) {
		case int64:
			cfg.ClientMaxRate = val
		case int:
			cfg.ClientMaxRate = int64(val)
		case float64:
			cfg.ClientMaxRate = int64(val)
		}
	}
	if v, ok := tmp["key"]; ok {
		cfg.Key = fmt.Sprintf("%v", v)
	}
	return cfg
}
