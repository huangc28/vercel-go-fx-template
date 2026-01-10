package cache

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/example/vercel-go-service-template/config"
)

func NewRedis(lc fx.Lifecycle, cfg config.Config, logger *zap.Logger) (*redis.Client, error) {
	if cfg.RedisURL == "" {
		logger.Info("redis disabled (redis_url not set)")
		return nil, nil
	}

	opt, err := redis.ParseURL(cfg.RedisURL)
	if err != nil {
		return nil, err
	}

	opt.ReadTimeout = 2 * time.Second
	opt.WriteTimeout = 2 * time.Second
	opt.DialTimeout = 2 * time.Second
	opt.PoolSize = 5
	opt.MinIdleConns = 1

	client := redis.NewClient(opt)

	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			return client.Close()
		},
	})

	return client, nil
}
