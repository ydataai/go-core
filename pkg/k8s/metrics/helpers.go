package metrics

import (
	"regexp"

	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
)

func newGauge(name string, labels []string) *prometheus.GaugeVec {
	return prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: metricNameFormat(name),
			Help: "Custom Gauge as helpers",
		},
		labels,
	)
}

func newCounter(name string, labels []string) *prometheus.CounterVec {
	return prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: metricNameFormat(name),
			Help: "Custom Counter as helpers",
		},
		labels,
	)
}

func metricNameFormat(name string) string {
	re, err := regexp.Compile(`[^\w]`)
	if err != nil {
		log.Errorf("An unexpected error occurred while generating a metric name: %v", err)
		panic("please, insert another metric name")
	}

	return re.ReplaceAllString(name, "_")
}
