package clients

import "github.com/kelseyhightower/envconfig"

// RedisConfiguration represents the client configuration to connect to Prometheus.
type RedisConfiguration struct {
	// Address represents host:port list separated by ,
	Address    []string `envconfig:"REDIS_ADDRESS" required:"true"`
	MasterName string   `envconfig:"REDIS_MASTER_NAME"`
	CACert     string   `envconfig:"REDIS_CA_CERT"`
	Cert       string   `envconfig:"REDIS_CERT"`
	CertKey    string   `envconfig:"REDIS_CERT_KEY"`
}

// LoadFromEnvVars for RedisConfiguration.
func (c *RedisConfiguration) LoadFromEnvVars() error {
	return envconfig.Process("", c)
}
