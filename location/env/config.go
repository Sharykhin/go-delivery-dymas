package env

import (
	coreenv "github.com/caarlos0/env/v8"
)

// Config Provides values env.
type Config struct {
	KafkaAddress                                 string `env:"KAFKA_BROKERS" envDefault:"localhost:9092"`
	KafkaSchemaRegistryAddress                   string `env:"KAFKA_SCHEMA_REGISTRY_ADDRESS" envDefault:"http://localhost:8085"`
	CourierGrpcAddress                           string `env:"COURIER_GRPC_ADDRESS" envDefault:":9666"`
	RedisAddress                                 string `env:"REDIS_ADDRESS" envDefault:"localhost:6379"`
	PortServer                                   string `env:"PORT_SERVER" envDefault:"8081"`
	Db                                           int    `env:"DB_REDIS" envDefault:"0"`
	Assignor                                     string `env:"KAFKA_CONSUMER_ASSIGNOR" envDefault:"range"`
	Oldest                                       bool   `env:"KAFKA_CONSUMER_OLDEST" envDefault:"true"`
	Verbose                                      bool   `env:"KAFKA_CONSUMER_VERBOSE" envDefault:"false"`
	PostgresDB                                   string `env:"POSTGRES_DB" envDefault:"courier_location"`
	PostgresPassword                             string `env:"POSTGRES_PASSWORD" envDefault:"S3cret"`
	PostgresUser                                 string `env:"POSTGRES_USER" envDefault:"citizix_user"`
	CourierLocationQueueSizeTasks                int    `env:"COURIER_LOCATION_QUEUE_SIZE_TASKS" envDefault:"10000"`
	CourierLocationWorkerPoolCount               int    `env:"COURIER_LOCATION_WORKER_POOL_COUNT" envDefault:"10"`
	CourierLocationWorkerTimeoutGracefulShutdown int    `env:"COURIER_LOCATION_WORKER_TIMEOUT_GRACEFUL_SHUTDOWN" envDefault:"30"`
}

// GetConfig gets Configs env.
func GetConfig() (config Config, err error) {
	cfg := Config{}
	err = coreenv.Parse(&cfg)

	return cfg, err
}
