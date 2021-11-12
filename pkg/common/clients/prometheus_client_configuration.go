package clients

import "github.com/kelseyhightower/envconfig"

// PrometheusConfiguration represents the client configuration to connect to Prometheus.
type PrometheusConfiguration struct {
	Address string `envconfig:"PROMETHEUS_ADDRESS" required:"true"`
}

// LoadFromEnvVars for PrometheusConfiguration.
func (c *PrometheusConfiguration) LoadFromEnvVars() error {
	return envconfig.Process("", c)
}
