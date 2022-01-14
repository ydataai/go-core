package clients

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/ydataai/go-core/pkg/common/logging"
)

// RedisClient represents the Redis client.
type RedisClient struct {
	*redis.Client
}

// NewRedisClient creates a new RedisClient (redis.Client) instance.
func NewRedisClient(config RedisConfiguration, logger logging.Logger) RedisClient {
	ctx := context.Background()
	client := redis.NewClient(&redis.Options{
		Addr: config.Address,
	})
	err := client.Ping(ctx).Err()
	if err != nil {
		logger.Fatalf("Error while connect to Redis: %s", config.Address)
	}
	// make sure the redis server is ready to write.
	err = client.Set(ctx, "lastUpdate", time.Now(), time.Minute).Err()
	if err != nil {
		logger.Fatalf("Redis Server is ready-only. Err: %v", err)
	}
	return RedisClient{client}
}
