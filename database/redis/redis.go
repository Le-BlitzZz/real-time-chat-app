package redis

import (
	"context"

	go_redis "github.com/redis/go-redis/v9"
)

var nLastMessages int64 = 10

type RedisDb struct {
	*go_redis.Client
}

func New(address string) (*RedisDb, error) {
	redisDb := &RedisDb{go_redis.NewClient(&go_redis.Options{
		Addr: address,
	})}

	if err := redisDb.Ping(context.Background()).Err(); err != nil {
		return nil, err
	}

	return redisDb, nil
}
