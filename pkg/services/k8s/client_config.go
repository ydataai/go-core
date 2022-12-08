package k8s

import (
	"github.com/kelseyhightower/envconfig"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/util/flowcontrol"
	ctrl "sigs.k8s.io/controller-runtime"
)

// K8sRestClientConfiguration defines required variables to configure the environment
type K8sRestClientConfiguration struct {
	K8S_CLIENT_QPS   float32 `envconfig:"K8S_REST_CLIENT_QPS" default:"100"`
	K8S_CLIENT_BURST int     `envconfig:"K8S_REST_CLIENT_BURST" default:"500"`
}

// LoadFromEnvVars for K8sRestClientConfiguration
func (c *K8sRestClientConfiguration) LoadFromEnvVars() error {
	return envconfig.Process("", c)
}

func Config(config K8sRestClientConfiguration) *rest.Config {
	kconfig := ctrl.GetConfigOrDie()
	kconfig.RateLimiter = flowcontrol.NewTokenBucketRateLimiter(config.K8S_CLIENT_QPS, config.K8S_CLIENT_BURST)
	return kconfig
}
