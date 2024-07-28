package env

import (
	coreenv "github.com/caarlos0/env/v8"
)

// Config Provides values env.
type Config struct {
	KafkaAddress               string `env:"KAFKA_BROKERS" envDefault:"localhost:9092"`
	KafkaSchemaRegistryAddress string `env:"KAFKA_SCHEMA_REGISTRY_ADDRESS" envDefault:"http://localhost:8085"`
	CourierGrpcAddress         string `env:"COURIER_GRPC_ADDRESS" envDefault:":9666"`
}

// GetConfig gets Configs env.
func GetConfig() (config Config, err error) {
	cfg := Config{}
	err = coreenv.Parse(&cfg)

	return cfg, err
}
