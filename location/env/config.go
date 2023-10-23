package env

import (
	coreenv "github.com/caarlos0/env/v8"
)

type Config struct {
	KafkaAddress       string `env:"KAFKA_BROKERS" envDefault:"localhost:9092"`
	CourierGrpcAddress string `env:"COURIER_GRPC_ADDRESS" envDefault:":9666"`
	RedisAddress       string `env:"REDIS_ADDRESS" envDefault:"localhost:6379"`
	PortServer         string `env:"PORT_SERVER" envDefault:"8081"`
	Db                 int    `env:"DB_REDIS" envDefault:"0"`
	Assignor           string `env:"KAFKA_CONSUMER_ASSIGNOR" envDefault:"range"`
	Oldest             bool   `env:"KAFKA_CONSUMER_OLDEST" envDefault:"true"`
	Verbose            bool   `env:"KAFKA_CONSUMER_VERBOSE" envDefault:"false"`
	DbName             string `env:"POSTGRES_DB" envDefault:"courier_location"`
	DbPassword         string `env:"POSTGRES_PASSWORD" envDefault:"S3cret"`
	DbUser             string `env:"POSTGRES_USER" envDefault:"citizix_user"`
}

func GetConfig() (config Config, err error) {
	cfg := Config{}
	err = coreenv.Parse(&cfg)

	return cfg, err
}
