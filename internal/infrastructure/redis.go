package infrastructure

import (
	"context"

	"github.com/redis/go-redis/v9"
)

func NewRedisClient(addr, pass string) (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: pass,
		DB:       0,
	})

	if _, err := rdb.Ping(context.Background()).Result(); err != nil {
		return nil, err
	}
	return rdb, nil
}
