package redis

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/ydataai/go-core/pkg/common/logging"
)

type RedisClient interface {
	Get(ctx context.Context, key string) *redis.StringCmd
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd
	Publish(ctx context.Context, channel string, message interface{}) *redis.IntCmd
	Subscribe(ctx context.Context, channels ...string) *redis.PubSub
	Ping(ctx context.Context) *redis.StatusCmd
}

// RedisClient represents the Redis client.
type redisClientImpl struct {
	client  *redis.Client
	cluster *redis.ClusterClient
}

// NewRedisClient creates a new RedisClient (redis.Client) instance.
func NewRedisClient(config RedisConfiguration, logger logging.Logger) RedisClient {
	var client RedisClient
	var err error

	if config.CACert != "" && config.Cert != "" && config.CertKey != "" {
		client = newRedisClusterClient(config, logger)
	} else {
		client = newRedisSingleNodeClient(config)
	}

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

	return client
}

func newRedisClusterClient(config RedisConfiguration, logger logging.Logger) RedisClient {
	// CA and Cert configuration for TLS connection
	certPool := x509.NewCertPool()
	caCert, err := os.ReadFile(config.CACert)
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
	return redisClientImpl{
		cluster: redis.NewClusterClient(&redis.ClusterOptions{
			Addrs: config.Address,
			TLSConfig: &tls.Config{
				MinVersion: tls.VersionTLS12,
				RootCAs:    certPool,
				Certificates: []tls.Certificate{
					cert,
				},
				InsecureSkipVerify: config.InsecureSkipVerify,
			},
		}),
	}
}

func newRedisSingleNodeClient(config RedisConfiguration) RedisClient {
	return redisClientImpl{
		client: redis.NewClient(&redis.Options{
			Addr: config.Address[0],
		})}
}

func (c redisClientImpl) get() RedisClient {
	if c.cluster != nil {
		return c.cluster
	}
	return c.client
}

func (c redisClientImpl) Get(ctx context.Context, key string) *redis.StringCmd {
	return c.get().Get(ctx, key)
}

func (c redisClientImpl) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd {
	return c.get().Set(ctx, key, value, expiration)
}

func (c redisClientImpl) Publish(ctx context.Context, channel string, message interface{}) *redis.IntCmd {
	return c.get().Publish(ctx, channel, message)
}

func (c redisClientImpl) Subscribe(ctx context.Context, channels ...string) *redis.PubSub {
	return c.get().Subscribe(ctx, channels...)
}

func (c redisClientImpl) Ping(ctx context.Context) *redis.StatusCmd {
	return c.get().Ping(ctx)
}
