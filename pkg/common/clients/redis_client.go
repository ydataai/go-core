package clients

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/ydataai/go-core/pkg/common/logging"
)

type RedisClient interface {
	Get(ctx context.Context, key string) *redis.StringCmd
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd
	Publish(ctx context.Context, channel string, message interface{}) *redis.IntCmd
	Subscribe(ctx context.Context, channels ...string) *redis.PubSub
}

// RedisClient represents the Redis client.
type redisClientImpl struct {
	client *redis.Client
}

// NewRedisClient creates a new RedisClient (redis.Client) instance.
func NewRedisClient(config RedisConfiguration, logger logging.Logger) RedisClient {
	// CA and Cert configuration for TLS connection
	certPool := x509.NewCertPool()
	caCert, err := ioutil.ReadFile(config.CACert)
	if err != nil {
		logger.Fatalf("Error to read CA Cert file from: %s. Err: %v", config.CACert, err)
	}
	if ok := certPool.AppendCertsFromPEM(caCert); !ok {
		logger.Fatalf("Error to add CA Cert to cert pool. Err: %v", err)
	}
	cert, err := tls.LoadX509KeyPair(config.Cert, config.CertKey)
	if err != nil {
		logger.Fatalf("Error to read Redis Cert file from: %s, %s. Err: %v", config.Cert, config.CertKey, err)
	}
	// redis client initialization
	client := redis.NewFailoverClient(&redis.FailoverOptions{
		MasterName:    config.MasterName,
		SentinelAddrs: config.Address,
		TLSConfig: &tls.Config{
			MinVersion: tls.VersionTLS12,
			RootCAs:    certPool,
			Certificates: []tls.Certificate{
				cert,
			},
		},
	})
	ctx := context.Background()
	// test server with ping/pong
	err = client.Ping(ctx).Err()
	if err != nil {
		logger.Fatalf("Error while connect to Redis: %s. Err: %v", config.Address, err)
	}
	// make sure the redis server is ready to write.
	err = client.Set(ctx, "lastUpdate", time.Now(), time.Minute).Err()
	if err != nil {
		logger.Fatalf("Redis Server is ready-only. Err: %v", err)
	}
	return redisClientImpl{client: client}
}

func (c redisClientImpl) Get(ctx context.Context, key string) *redis.StringCmd {
	return c.client.Get(ctx, key)
}

func (c redisClientImpl) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd {
	return c.client.Set(ctx, key, value, expiration)
}

func (c redisClientImpl) Publish(ctx context.Context, channel string, message interface{}) *redis.IntCmd {
	return c.client.Publish(ctx, channel, message)
}

func (c redisClientImpl) Subscribe(ctx context.Context, channels ...string) *redis.PubSub {
	return c.client.Subscribe(ctx, channels...)
}
