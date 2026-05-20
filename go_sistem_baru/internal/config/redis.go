package config

import (
	"context"

	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
)

func NewRedis(cfg *Config, log *logrus.Logger) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.Addr,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
		PoolSize: cfg.Redis.PoolSize,
	})

	if err := rdb.Ping(context.Background()).Err(); err != nil {
		log.WithError(err).Fatal("failed to connect to redis")
	}

	log.Info("redis connection established")
	return rdb
}
