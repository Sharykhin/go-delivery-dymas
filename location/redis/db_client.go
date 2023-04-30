package redis

import (
	coreredis "github.com/redis/go-redis/v9"
)

var Config = &coreredis.Options{
	Addr: "localhost:6379",
	DB:   0,
}

func CreateConnect(options *coreredis.Options) *coreredis.Client {
	// create a new Redis client
	client := coreredis.NewClient(options)

	return client
}
