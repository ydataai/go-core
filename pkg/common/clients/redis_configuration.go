package clients

import "github.com/kelseyhightower/envconfig"

// RedisConfiguration represents the client configuration to connect to Prometheus.
type RedisConfiguration struct {
	// Address represents host:port
	Address string `envconfig:"REDIS_ADDRESS" required:"true"`
}

// LoadFromEnvVars for RedisConfiguration.
func (c *RedisConfiguration) LoadFromEnvVars() error {
	return envconfig.Process("", c)
}
