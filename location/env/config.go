package env

import (
	coreenv "github.com/caarlos0/env/v8"
)

type Config struct {
	Address    string `env:"KAFKA_BROKERS" envDefault:"localhost:9092"`
	Addr       string `env:"REDIS_ADDRESS" envDefault:"localhost:6379"`
	PortServer string `env:"PORT_SERVER" envDefault:"8081"`
	Assignor   string `env:"ASSIGNOR" envDefault:"range"`
	Oldest     bool   `env:"OLDEST" envDefault:"true"`
	Verbose    bool   `env:"VERBOSE" envDefault:"false"`
	DbName     string `env:"POSTGRES_DB" envDefault:"courier_location"`
	PasswordDb string `env:"POSTGRES_PASSWORD" envDefault:"S3cret"`
	DbUser     string `env:"POSTGRES_USER" envDefault:"citizix_user"`
}

func GetConfig() (config Config, err error) {
	cfg := Config{}
	err = coreenv.Parse(&cfg)

	return cfg, err
}
