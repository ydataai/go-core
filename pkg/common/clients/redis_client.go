package clients

import (
	"context"

	"github.com/go-redis/redis/v8"
	"github.com/ydataai/go-core/pkg/common/logging"
)

// RedisClient represents the Redis client.
type RedisClient struct {
	*redis.Client
}

// NewRedisClient creates a new RedisClient (redis.Client) instance.
func NewRedisClient(config RedisConfiguration, logger logging.Logger) RedisClient {
	client := redis.NewClient(&redis.Options{
		Addr: config.Address,
	})
	err := client.Ping(context.Background()).Err()
	if err != nil {
		logger.Fatalf("Error while connect to Redis: %s", config.Address)
	}
	return RedisClient{client}
}
