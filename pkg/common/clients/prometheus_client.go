package clients

import (
	"context"
	"time"

	"github.com/prometheus/client_golang/api"
	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/common/model"
	"github.com/sirupsen/logrus"
)

// PrometheusClient represents the Prometheus client.
type PrometheusClient struct {
	logger *logrus.Logger
	config PrometheusConfiguration
	api    v1.API
}

// NewPrometheusClient creates a PrometheusClient instance.
func NewPrometheusClient(
	config PrometheusConfiguration,
	logger *logrus.Logger) PrometheusClient {
	api := createClientAPI(config, logger)
	return PrometheusClient{
		config: config,
		logger: logger,
		api:    api,
	}
}

func createClientAPI(config PrometheusConfiguration, logger *logrus.Logger) v1.API {
	client, err := api.NewClient(api.Config{
		Address: config.Address,
	})
	if err != nil {
		logger.Errorf("Error creating client: %v", err)
	}
	return v1.NewAPI(client)
}

// QueryRange returns range vectors as result type matrix, given the following parameters:
//    - query: Prometheus query
//    - startAt: start time
//    - endAt: end time
//    - step: interval duration
// It returns the result matrix and the execution error encountered.
func (c PrometheusClient) QueryRange(ctx context.Context, query string, startAt, endAt time.Time, step time.Duration) (model.Value, error) {
	result, warnings, err := c.api.QueryRange(ctx, query, v1.Range{
		Start: startAt,
		End:   endAt,
		Step:  step,
	})
	if err != nil {
		c.logger.Errorf("Error querying Prometheus: %v", err)
		return nil, err
	}
	if len(warnings) > 0 {
		c.logger.Warnf("Prometheus query warnings result: %v", warnings)
	}
	return result, nil
}