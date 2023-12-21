package env

import (
	coreenv "github.com/caarlos0/env/v8"
)

// Config provides values env
type Config struct {
	DBName           string `env:"POSTGRES_DB" envDefault:"orders"`
	DBPassword       string `env:"POSTGRES_PASSWORD" envDefault:"S3cret"`
	DBUser           string `env:"POSTGRES_USER" envDefault:"citizix_user"`
	PortServerOrder  string `env:"PORT_SERVER_ORDER" envDefault:"6661"`
	OrderGrpcAddress string `env:"ORDER_GRPC_ADDRESS" envDefault:":9667"`
	KafkaAddress     string `env:"KAFKA_BROKERS" envDefault:"localhost:9092"`
	Assignor         string `env:"KAFKA_CONSUMER_ASSIGNOR" envDefault:"range"`
	Oldest           bool   `env:"KAFKA_CONSUMER_OLDEST" envDefault:"true"`
	Verbose          bool   `env:"KAFKA_CONSUMER_VERBOSE" envDefault:"false"`
}

// GetConfig gets Configs env
func GetConfig() (config Config, err error) {
	cfg := Config{}
	err = coreenv.Parse(&cfg)

	return cfg, err
}
