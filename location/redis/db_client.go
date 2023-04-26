package redis

import (
	"context"
	coreredis "github.com/redis/go-redis/v9"
)

func CreateConnect(ctx context.Context, options *coreredis.Options) (*coreredis.Client, error) {
	// create a new Redis client
	client := coreredis.NewClient(options)

	// use the PING command to check the connection
	err := client.Ping(ctx).Err()
	if err != nil {
		return client, err
	}

	return client, nil
}
