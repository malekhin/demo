package app

import (
	"demo/internal/common/config"
	redisPkg "demo/pkg/redis"
	"fmt"

	"github.com/redis/go-redis/v9"

	"go.uber.org/fx"
)

func NewRedis(cfg *config.Config) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%d", cfg.Redis.Host, cfg.Redis.Port),
	})
}

func NewRedisPkg(cfg *config.Config) redisPkg.IRedis {
	return redisPkg.NewRedis(redisPkg.Config{
		Host:     cfg.Redis.Host,
		Port:     cfg.Redis.Port,
		Username: cfg.Redis.Username,
		Password: cfg.Redis.Password,
		Prefix:   cfg.Redis.Prefix,
	})
}

var RedisModule = fx.Module("redis",
	fx.Provide(NewRedis),
	fx.Provide(NewRedisPkg),
)
