package env

import (
	coreenv "github.com/caarlos0/env/v8"
)

type Config struct {
	DbName            string `env:"POSTGRES_DB" envDefault:"couriers"`
	DbPassword        string `env:"POSTGRES_PASSWORD" envDefault:"S3cret"`
	DbUser            string `env:"POSTGRES_USER" envDefault:"citizix_user"`
	PortServerCourier string `env:"PORT_SERVER_COURIER" envDefault:"8881"`
}

func GetConfig() (config Config, err error) {
	cfg := Config{}
	err = coreenv.Parse(&cfg)

	return cfg, err
}
