package clients

import "github.com/kelseyhightower/envconfig"

// PrometheusConfiguration represents the client configuration to connect to Prometheus.
type PrometheusConfiguration struct {
	Address string `envconfig:"PROMETHEUS_ADDRESS" required:"true"`
}

func (c *PrometheusConfiguration) LoadFromEnvVars() error {
	if err := envconfig.Process("", c); err != nil {
		return err
	}
	return nil
}
