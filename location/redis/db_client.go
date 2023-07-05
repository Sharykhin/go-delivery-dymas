package redis

import (
	coreredis "github.com/redis/go-redis/v9"
)

func NewConnect(addr string, db int) *coreredis.Client {

	var options = &coreredis.Options{
		Addr: addr,
		DB:   db,
	}

	client := coreredis.NewClient(options)

	return client
}
