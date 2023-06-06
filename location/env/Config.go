package env

import (
	"os"
	"strconv"
)

func GetKafkaConfig() (address string) {
	address = "localhost:9092"
	if os.Getenv("KAFKA_BROKERS") != "" {
		address = os.Getenv("KAFKA_BROKERS")
	}

	return address
}

func GetRedisConfig() (addr string, db int) {
	host := "localhost"
	port := "6379"
	db = 0
	addr = host + ":" + port
	if os.Getenv("REDIS_HOST") != "" && os.Getenv("REDIS_PORT") != "" {
		host = os.Getenv("REDIS_HOST")
		port = os.Getenv("REDIS_PORT")
		addr = host + ":" + port
	}

	if os.Getenv("DB_REDIS") != "" {
		dbConvert, error := strconv.Atoi(os.Getenv("DB_REDIS"))
		db = dbConvert
		if error != nil {
			panic(error)
		}
	}

	return addr, db
}

func GetServerEnv() (portServer string) {
	portServer = ":8081"

	if os.Getenv("HTTP_PORT") != "" {
		portServer = os.Getenv("HTTP_PORT")
	}

	return portServer
}
