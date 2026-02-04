package redis

import (
	"agrigation_api/pkg/tools"
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

type Redis struct {
	pool *redis.Client
}

func NewRedis() (*Redis, error) {
	rdsHost := tools.GetEnv("REDIS_HOST", "localhost")
	rdsPort := tools.GetEnvAsInt("REDIS_PORT", 6379)
	rdsPassword := tools.GetEnv("REDIS_PASSWORD", "")
	rdsDB := tools.GetEnvAsInt("REDIS_DB", 0)
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", rdsHost, rdsPort),
		Password: rdsPassword,
		DB:       rdsDB,
	})
	ctx := context.Background()
	_, err := client.Ping(ctx).Result()
	if err != nil {
		return nil, err
	}
	return &Redis{
		pool: client,
	}, nil
}
