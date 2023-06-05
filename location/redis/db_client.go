package redis

import (
	coreredis "github.com/redis/go-redis/v9"
)

func CreateConnect(addr string, db int) *coreredis.Client {

	var options = &coreredis.Options{
		Addr: addr,
		DB:   db,
	}
	// create a new Redis client
	client := coreredis.NewClient(options)

	return client
}
