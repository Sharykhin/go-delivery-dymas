package env

import (
	coreenv "github.com/caarlos0/env/v8"
)

// Config Provides values env.
type Config struct {
	PortServerGateway string `env:"PORT_SERVER_GATEWAY" envDefault:"8080"`
}

// GetConfig gets Configs env.
func GetConfig() (config Config, err error) {
	cfg := Config{}
	err = coreenv.Parse(&cfg)

	return cfg, err
}
