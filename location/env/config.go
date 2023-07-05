package env

import (
	coreenv "github.com/caarlos0/env/v8"
)

type Config struct {
	Address    string `env:"KAFKA_BROKERS" envDefault:"localhost:9092"`
	Addr       string `env:"REDIS_ADDRESS" envDefault:"localhost:6379"`
	PortServer string `env:"PORT_SERVER" envDefault:"8081"`
	Db         int    `env:"DB_REDIS" envDefault:"0"`
}

func GetConfig() (config Config, err error) {
	cfg := Config{}
	err = coreenv.Parse(&cfg)

	return cfg, err
}
