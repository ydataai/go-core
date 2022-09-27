package config

import (
	"github.com/kelseyhightower/envconfig"
)

// ManagerConfiguration defines required variables to configure the environment
type ManagerConfiguration struct {
	EnableLeaderElection bool   `envconfig:"ENABLE_LEADER_ELECTION" required:"true"`
	LeaderElectionID     string `envconfig:"LEADER_ELECTION_ID" required:"true"`
	Port                 int    `envconfig:"MANAGER_PORT" default:"9443"`
	MetricsServerPort    int    `envconfig:"METRICS_SERVER_PORT" default:"8080"`
	HealthProbeAddress   string `envconfig:"HEALTH_PROBE_ADDRESS" default:":8081"`
	EnableWebhooks       bool   `envconfig:"ENABLE_WEBHOOKS" default:"true"`
	WebhookCertDir       string `envconfig:"WEBHOOK_CERT_DIR" default:"/tmp/k8s-webhook-server/serving-certs"`
}

// LoadFromEnvVars from the Manager
func (c *ManagerConfiguration) LoadFromEnvVars() error {
	return envconfig.Process("", c)
}
