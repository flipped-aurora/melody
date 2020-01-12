package config

//Parser returns a ServiceConfig struct according to the providing configFile name.
type Parser interface {
	Parse(configFile string) (ServiceConfig, error)
}
