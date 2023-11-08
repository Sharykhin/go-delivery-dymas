package env

import (
	coreenv "github.com/caarlos0/env/v8"
)

// Config Provides values env
type Config struct {
	DBName             string `env:"POSTGRES_DB" envDefault:"courier"`
	DBPassword         string `env:"POSTGRES_PASSWORD" envDefault:"S3cret"`
	DBUser             string `env:"POSTGRES_USER" envDefault:"citizix_user"`
	PortServerCourier  string `env:"PORT_SERVER_COURIER" envDefault:"8881"`
	CourierGrpcAddress string `env:"COURIER_GRPC_ADDRESS" envDefault:":9666"`
}

// GetConfig gets Configs env
func GetConfig() (config Config, err error) {
	cfg := Config{}
	err = coreenv.Parse(&cfg)

	return cfg, err
}
