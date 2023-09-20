package flags

import (
	"context"

	"github.com/redis/go-redis/v9"
)

type RedisFlags struct {
	Addr     string `toml:"addr"`
	Password string `toml:"password"`
}

func (f RedisFlags) Init() (*redis.Client, error) {
	cfg := redis.Options{Addr: f.Addr, Password: f.Password, DB: 0}
	redisClient := redis.NewClient(&cfg)

	err := redisClient.Ping(context.TODO()).Err()

	if err != nil {
		return nil, err
	}
	return redisClient, nil
}
