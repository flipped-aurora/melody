package config

import "github.com/spf13/viper"

func New() ViperParser {
	return ViperParser{
		viper: viper.New(),
	}
}

type ViperParser struct {
	viper *viper.Viper
}

func (p ViperParser) Parse(configFile string) (ServiceConfig, error) {
	return ServiceConfig{}, nil
}
