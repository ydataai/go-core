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

// RedisClient represents the Redis client.
type RedisClient struct {
	*redis.Client
}

// NewRedisClient creates a new RedisClient (redis.Client) instance.
func NewRedisClient(config RedisConfiguration, logger logging.Logger) RedisClient {
	// CA cert configuration for TLS connection
	certPool := x509.NewCertPool()
	caCert, err := ioutil.ReadFile(config.CACert)
	if err != nil {
		logger.Fatalf("Error to read CA Cert file from: %s. Err: %v", config.CACert, err)
	}
	if ok := certPool.AppendCertsFromPEM(caCert); !ok {
		logger.Fatalf("Error to add CA Cert to cert pool. Err: %v", err)
	}
	// redis client initialization
	client := redis.NewClient(&redis.Options{
		Addr: config.Address,
		TLSConfig: &tls.Config{
			MinVersion: tls.VersionTLS12,
			RootCAs:    certPool,
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
	return RedisClient{client}
}
