package config

import (
	"github.com/spf13/viper"
	"reflect"
	"unsafe"
)

//New return new parser extends viper
func New() ViperParser {
	return ViperParser{
		viper: viper.New(),
	}
}

//ViperParser extends viper
type ViperParser struct {
	viper *viper.Viper
}



//Parse to parse config file
func (p ViperParser) Parse(configFile string) (ServiceConfig, error) {
	p.viper.SetConfigFile(configFile)
	p.viper.AutomaticEnv()
	var cfg ServiceConfig

	if err := p.viper.ReadInConfig(); err != nil {
		return cfg, checkErr(err, configFile)
	}

	if err := p.viper.Unmarshal(&cfg); err != nil {
		return cfg, checkErr(err, configFile)
	}

	if err := cfg.Init(); err != nil {
		return cfg, CheckErr(err, configFile)
	}

	return cfg, nil
}

func checkErr(err error, configFile string) error {
	switch e := err.(type) {
	case viper.ConfigParseError:
		var subErr error
		rs := reflect.ValueOf(&e).Elem()
		rf := rs.Field(0)
		ri := reflect.ValueOf(&subErr).Elem()

		rf = reflect.NewAt(rf.Type(), unsafe.Pointer(rf.UnsafeAddr())).Elem()

		ri.Set(rf)

		return checkErr(subErr, configFile)
	default:
		return CheckErr(err, configFile)
	}
}