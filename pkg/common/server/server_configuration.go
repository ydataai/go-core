package server

import (
	"time"

	"github.com/kelseyhightower/envconfig"
)

// HTTPServerConfiguration is a struct that holds all the environment variables required to the HTTP server
type HTTPServerConfiguration struct {
	Host                string        `envconfig:"HTTP_HOST" default:""`
	Port                int           `envconfig:"HTTP_PORT" default:"80"`
	RequestTimeout      time.Duration `envconfig:"HTTP_REQUEST_TIMEOUT" default:"30s"`
	UserID              string        `envconfig:"USER_ID" default:""`
	CertificateFile     string        `envconfig:"CERTIFICATE_FILE" default:"/etc/tls/tls.crt"`
	CertificateKeyFile  string        `envconfig:"CERTIFICATE_KEY_FILE" default:"/etc/tls/tls.key"`
	HealthCheckEndpoint string        `envconfig:"HEALTH_CHECK_ENDPOINT" default:"/healthz"`
}

// LoadFromEnvVars reads all env vars required for the server package
func (h *HTTPServerConfiguration) LoadFromEnvVars() error {
	return envconfig.Process("", h)
}
