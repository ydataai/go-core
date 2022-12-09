package k8s

import (
	"github.com/kelseyhightower/envconfig"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/util/flowcontrol"
	ctrl "sigs.k8s.io/controller-runtime"
)

// RestClientConfiguration defines required variables to configure the environment
type RestClientConfiguration struct {
	ClientQPS   float32 `envconfig:"K8S_REST_CLIENT_QPS" default:"100"`
	ClientBurst int     `envconfig:"K8S_REST_CLIENT_BURST" default:"500"`
}

// LoadFromEnvVars for RestClientConfiguration.
func (c *RestClientConfiguration) LoadFromEnvVars() error {
	return envconfig.Process("", c)
}

// Config creates rest client configuration with TokenBucketRateLimiter.
func Config(config RestClientConfiguration) *rest.Config {
	kconfig := ctrl.GetConfigOrDie()
	kconfig.RateLimiter = flowcontrol.NewTokenBucketRateLimiter(config.ClientQPS, config.ClientBurst)
	return kconfig
}
