package config

const (
	BracketsRouterPatternBuilder = iota
	ColonRouterPatternBuilder
)

//ServiceConfig contains all config in melody server.
type ServiceConfig struct {
	ExtraConfig ExtraConfig `mapstructure:"extra_config"`
	Port        int         `mapstructure:"port"`
	//melody is in debug model
	Debug bool
}

//Extra config for melody
type ExtraConfig map[string]interface{}

func (s *ServiceConfig) Init() error {
	//TODO 初始化URIParser
	//TODO 判断版本一致
	//TODO 初始化全局参数
	//TODO 初始化Endpoints
	return nil
}
