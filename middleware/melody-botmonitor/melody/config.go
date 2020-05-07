package melody

import (
	"encoding/json"
	"errors"

	"melody/config"
	botmonitor "melody/middleware/melody-botmonitor"
)

// Namespace 命名空间
const Namespace = "melody_botmonitor"

// ErrNoConfig 当没有为模块定义配置时，将返回ErrNoConfig
var ErrNoConfig = errors.New("no config defined for the module")

// ParseConfig 从ExtraConfig 中提取模块配置，并返回一个适合使用botmonitor包的结构
func ParseConfig(cfg config.ExtraConfig) (botmonitor.Config, error) {
	res := botmonitor.Config{}
	e, ok := cfg[Namespace]
	if !ok {
		return res, ErrNoConfig
	}
	b, err := json.Marshal(e)
	if err != nil {
		return res, err
	}
	err = json.Unmarshal(b, &res)
	return res, err
}
